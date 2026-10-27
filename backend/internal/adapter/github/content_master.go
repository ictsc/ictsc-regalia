package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"gopkg.in/yaml.v3"
)

// The content repository follows the ICTSC master layout:
//
//	schedule.yaml                  (path configurable)
//	rule.md                        (optional)
//	problems/<name>.problem.md     (frontmatter + Markdown body)
//	problems/<name>.explanation.md (optional)
//	notices/<slug>.md              (frontmatter + Markdown; _template.md is skipped)
//	images/<file>                  (referenced from Markdown as /images/<file>)
const (
	ruleFilePath          = "rule.md"
	problemsDir           = "problems"
	noticesDir            = "notices"
	imagesDir             = "images"
	problemFileSuffix     = ".problem.md"
	explanationFileSuffix = ".explanation.md"
	noticeTemplateSlug    = "_template"
	// exampleProblemCode marks the sample problem kept in the master
	// repository; it is never published.
	exampleProblemCode = "0000"
)

var imageNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,255}$`)

var imageContentTypes = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".webp": "image/webp",
}

type treeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	SHA  string `json:"sha"`
}

// MasterFiles holds the text files read from the content repository, keyed by
// repository path. Images are fetched separately.
type MasterFiles struct {
	Schedule     []byte
	Rule         []byte
	HasRule      bool
	Problems     map[string][]byte
	Explanations map[string][]byte
	Notices      map[string][]byte
}

func (c *Client) Fetch(ctx context.Context, repository, ref, commit string) (core.ContentSnapshot, error) {
	if repository != c.repository || ref != c.ref {
		return core.ContentSnapshot{}, errors.New("GitHub repository or ref is outside the configured content source")
	}
	owner, name, err := splitRepository(repository)
	if err != nil {
		return core.ContentSnapshot{}, err
	}
	if !commitPattern.MatchString(commit) {
		return core.ContentSnapshot{}, errors.New("GitHub commit must be a full lowercase object ID")
	}
	if err := c.verifyCommit(ctx, owner, name, commit); err != nil {
		return core.ContentSnapshot{}, err
	}
	tree, err := c.listTree(ctx, owner, name, commit)
	if err != nil {
		return core.ContentSnapshot{}, err
	}
	blobs := make(map[string]treeEntry, len(tree))
	for _, entry := range tree {
		if entry.Type == "blob" {
			blobs[entry.Path] = entry
		}
	}
	load := func(filePath string, maximum int64) ([]byte, error) {
		content, err := c.fetchFile(ctx, owner, name, commit, filePath, maximum)
		if err != nil {
			return nil, wrapContentLoadError(filePath, err)
		}
		return content, nil
	}
	files := MasterFiles{Problems: map[string][]byte{}, Explanations: map[string][]byte{}, Notices: map[string][]byte{}}
	if _, ok := blobs[c.schedulePath]; !ok {
		return core.ContentSnapshot{}, invalidManifest("schedule file %q is missing", c.schedulePath)
	}
	if files.Schedule, err = load(c.schedulePath, maxScheduleBytes); err != nil {
		return core.ContentSnapshot{}, err
	}
	if _, ok := blobs[ruleFilePath]; ok {
		if files.Rule, err = load(ruleFilePath, maxMarkdownBytes); err != nil {
			return core.ContentSnapshot{}, err
		}
		files.HasRule = true
	}
	var images []treeEntry
	for _, filePath := range sortedKeys(blobs) {
		if strings.HasPrefix(filePath, imagesDir+"/") {
			if err := validateImageEntry(blobs[filePath]); err != nil {
				return core.ContentSnapshot{}, err
			}
			images = append(images, blobs[filePath])
			continue
		}
		directory, base := path.Split(filePath)
		switch strings.TrimSuffix(directory, "/") {
		case problemsDir:
			switch {
			case strings.HasSuffix(base, problemFileSuffix):
				if files.Problems[filePath], err = load(filePath, maxMarkdownBytes); err != nil {
					return core.ContentSnapshot{}, err
				}
			case strings.HasSuffix(base, explanationFileSuffix):
				if files.Explanations[filePath], err = load(filePath, maxMarkdownBytes); err != nil {
					return core.ContentSnapshot{}, err
				}
			}
		case noticesDir:
			if strings.HasSuffix(base, ".md") && strings.TrimSuffix(base, ".md") != noticeTemplateSlug {
				if files.Notices[filePath], err = load(filePath, maxMarkdownBytes); err != nil {
					return core.ContentSnapshot{}, err
				}
			}
		}
	}
	manifest, err := BuildManifest(files, c.secondaryQualifier)
	if err != nil {
		return core.ContentSnapshot{}, err
	}
	assets, err := c.fetchImages(ctx, owner, name, images)
	if err != nil {
		return core.ContentSnapshot{}, err
	}
	return core.ContentSnapshot{
		CommitSHA:  commit,
		Repository: repository,
		Ref:        ref,
		Manifest:   manifest,
		FetchedAt:  c.now().UTC(),
		Assets:     assets,
	}, nil
}

func (c *Client) listTree(ctx context.Context, owner, repository, commit string) ([]treeEntry, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", c.apiBaseURL,
		url.PathEscape(owner), url.PathEscape(repository), url.PathEscape(commit))
	var response struct {
		Truncated bool        `json:"truncated"`
		Tree      []treeEntry `json:"tree"`
	}
	if err := c.getGitHubJSON(ctx, endpoint, "get repository tree", &response); err != nil {
		return nil, err
	}
	if response.Truncated {
		return nil, invalidManifest("repository tree is too large to list")
	}
	return response.Tree, nil
}

// validateImageEntry rejects anything in images/ that cannot be served, so a
// misplaced or unsupported file fails the refresh instead of a broken image.
func validateImageEntry(entry treeEntry) error {
	name := strings.TrimPrefix(entry.Path, imagesDir+"/")
	if strings.Contains(name, "/") {
		return invalidManifest("image %q must be directly under %s/", entry.Path, imagesDir)
	}
	if !imageNamePattern.MatchString(name) {
		return invalidManifest("image %q has a file name outside [A-Za-z0-9._-]", entry.Path)
	}
	if _, ok := imageContentTypes[strings.ToLower(path.Ext(name))]; !ok {
		return invalidManifest("image %q is not png, jpg, jpeg, gif or webp", entry.Path)
	}
	if !commitPattern.MatchString(entry.SHA) {
		return invalidManifest("image %q has an invalid blob SHA", entry.Path)
	}
	return nil
}

func (c *Client) fetchImages(ctx context.Context, owner, repository string, entries []treeEntry) ([]core.ContentAsset, error) {
	assets := make([]core.ContentAsset, 0, len(entries))
	var total int64
	for _, entry := range entries {
		if entry.Size > maxImageBytes {
			return nil, invalidManifest("image %q exceeds %d bytes", entry.Path, maxImageBytes)
		}
		total += entry.Size
		if total > maxTotalImageBytes {
			return nil, invalidManifest("images exceed %d bytes in total", maxTotalImageBytes)
		}
		endpoint := fmt.Sprintf("%s/repos/%s/%s/git/blobs/%s", c.apiBaseURL,
			url.PathEscape(owner), url.PathEscape(repository), url.PathEscape(entry.SHA))
		var response struct {
			Encoding string `json:"encoding"`
			Content  string `json:"content"`
			Size     int64  `json:"size"`
		}
		if err := c.getGitHubJSON(ctx, endpoint, "get image blob", &response); err != nil {
			return nil, wrapContentLoadError(entry.Path, err)
		}
		if response.Encoding != "base64" {
			return nil, invalidManifest("image %q blob is not base64 encoded", entry.Path)
		}
		data, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(response.Content, "\n", ""))
		if err != nil || int64(len(data)) != entry.Size {
			return nil, invalidManifest("image %q blob is invalid", entry.Path)
		}
		name := path.Base(entry.Path)
		assets = append(assets, core.ContentAsset{
			Name: name, ContentType: imageContentTypes[strings.ToLower(path.Ext(name))], SHA: entry.SHA, Data: data,
		})
	}
	return assets, nil
}

type scheduleFile struct {
	RankingFreezeAt *time.Time      `yaml:"ranking_freeze_at"`
	Schedules       []scheduleEntry `yaml:"schedules"`
}

type scheduleEntry struct {
	Name    string    `yaml:"name"`
	StartAt time.Time `yaml:"startAt"`
	EndAt   time.Time `yaml:"endAt"`
}

type problemFrontmatter struct {
	Code                    string                   `yaml:"code"`
	Title                   string                   `yaml:"title"`
	MaxScore                int32                    `yaml:"max_score"`
	Category                string                   `yaml:"category"`
	RedeployRule            *redeployRuleFrontmatter `yaml:"redeploy_rule"`
	SubmissionableSchedules []string                 `yaml:"submissionable_schedules"`
	Draft                   bool                     `yaml:"draft"`
	DefaultGraderDiscordID  string                   `yaml:"default_grader_discord_id"`
	AnswerMaxLength         *int32                   `yaml:"answer_max_length"`
	AllowAttachments        bool                     `yaml:"allow_attachments"`
}

type redeployRuleFrontmatter struct {
	Type              string `yaml:"type"`
	PenaltyThreshold  *int32 `yaml:"penalty_threshold"`
	PenaltyPercentage *int32 `yaml:"penalty_percentage"`
}

type noticeFrontmatter struct {
	Title string `yaml:"title"`
	From  string `yaml:"from"`
}

// BuildManifest converts the master repository files into a manifest.
// Secondary qualifier mode runs a single schedule that every problem belongs
// to, so problems may not name schedules and the ranking freeze is unused.
func BuildManifest(files MasterFiles, secondaryQualifier bool) (core.Manifest, error) {
	var schedule scheduleFile
	if err := decodeStrictYAML(files.Schedule, &schedule); err != nil {
		return core.Manifest{}, invalidManifest("schedule: %v", err)
	}
	sections := make([]core.Section, 0, len(schedule.Schedules))
	for _, entry := range schedule.Schedules {
		sections = append(sections, core.Section{Slug: entry.Name, Beginning: entry.StartAt.UTC(), Ending: entry.EndAt.UTC()})
	}
	if err := validateSchedule(sections); err != nil {
		return core.Manifest{}, err
	}
	if secondaryQualifier {
		if len(sections) != 1 {
			return core.Manifest{}, invalidManifest("secondary qualifier mode requires exactly one schedule, got %d", len(sections))
		}
		if schedule.RankingFreezeAt != nil {
			return core.Manifest{}, invalidManifest("ranking_freeze_at is not used in secondary qualifier mode")
		}
	}
	manifest := core.Manifest{Version: 1, Sections: sections, Problems: []core.Problem{}, Announcements: []core.Announcement{}}
	if schedule.RankingFreezeAt != nil {
		freezeAt := schedule.RankingFreezeAt.UTC()
		manifest.RankingFreezeAt = &freezeAt
	}
	if files.HasRule {
		manifest.RuleMarkdown = string(files.Rule)
		manifest.RulePath = ruleFilePath
	}
	sectionIndex := make(map[string]int, len(sections))
	for index, section := range sections {
		sectionIndex[section.Slug] = index
	}
	seenCodes := make(map[string]string)
	for _, filePath := range sortedKeys(files.Problems) {
		problem, publish, err := buildProblem(filePath, files.Problems[filePath], files.Explanations)
		if err != nil {
			return core.Manifest{}, err
		}
		if !publish {
			continue
		}
		if previous, exists := seenCodes[problem.value.Code]; exists {
			return core.Manifest{}, invalidManifest("problem code %q is used by %s and %s", problem.value.Code, previous, filePath)
		}
		seenCodes[problem.value.Code] = filePath
		slug, err := problemSchedule(problem, sections, secondaryQualifier)
		if err != nil {
			return core.Manifest{}, err
		}
		index, ok := sectionIndex[slug]
		if !ok {
			return core.Manifest{}, invalidManifest("problem %q references unknown schedule %q", problem.value.Code, slug)
		}
		problem.value.SectionSlug = slug
		manifest.Sections[index].ProblemIDs = append(manifest.Sections[index].ProblemIDs, problem.value.Code)
		manifest.Problems = append(manifest.Problems, problem.value)
	}
	for _, filePath := range sortedKeys(files.Notices) {
		announcement, err := buildNotice(filePath, files.Notices[filePath])
		if err != nil {
			return core.Manifest{}, err
		}
		manifest.Announcements = append(manifest.Announcements, announcement)
	}
	return manifest, nil
}

type parsedProblem struct {
	value     core.Problem
	schedules []string
}

func buildProblem(filePath string, data []byte, explanations map[string][]byte) (parsedProblem, bool, error) {
	var matter problemFrontmatter
	body, err := splitFrontmatter(data, &matter)
	if err != nil {
		return parsedProblem{}, false, invalidManifest("%s: %v", filePath, err)
	}
	code := strings.ToUpper(strings.TrimSpace(matter.Code))
	if matter.Draft || code == exampleProblemCode {
		return parsedProblem{}, false, nil
	}
	if !validIdentifier(code, 8) {
		return parsedProblem{}, false, invalidManifest("%s: code %q is invalid", filePath, matter.Code)
	}
	if strings.TrimSpace(matter.Title) == "" || !validTextLength(matter.Title, 1, 255) || !validTextLength(matter.Category, 0, 255) {
		return parsedProblem{}, false, invalidManifest("%s: title or category length is invalid", filePath)
	}
	if matter.MaxScore <= 0 {
		return parsedProblem{}, false, invalidManifest("%s: max_score must be positive", filePath)
	}
	if matter.DefaultGraderDiscordID != "" && !discordIDPattern.MatchString(matter.DefaultGraderDiscordID) {
		return parsedProblem{}, false, invalidManifest("%s: default_grader_discord_id is invalid", filePath)
	}
	if strings.TrimSpace(body) == "" {
		return parsedProblem{}, false, invalidManifest("%s: body is empty", filePath)
	}
	answerMaxLength := int32(core.DefaultAnswerMaxLength)
	if matter.AnswerMaxLength != nil {
		answerMaxLength = *matter.AnswerMaxLength
	}
	if answerMaxLength < 1 || answerMaxLength > core.MaxAnswerLength {
		return parsedProblem{}, false, invalidManifest("%s: answer_max_length must be between 1 and %d", filePath, core.MaxAnswerLength)
	}
	rule, err := redeployRule(matter.RedeployRule)
	if err != nil {
		return parsedProblem{}, false, invalidManifest("%s: redeploy_rule: %v", filePath, err)
	}
	problem := core.Problem{
		Code: code, Title: matter.Title, MaxScore: matter.MaxScore, Category: matter.Category,
		Type: "DESCRIPTIVE", Body: body, BodyPath: filePath, Redeploy: rule,
		DefaultGraderDiscordID: matter.DefaultGraderDiscordID,
		AnswerMaxLength:        answerMaxLength, AllowAttachments: matter.AllowAttachments,
	}
	explanationPath := strings.TrimSuffix(filePath, problemFileSuffix) + explanationFileSuffix
	if explanation, ok := explanations[explanationPath]; ok {
		problem.Explanation = string(explanation)
		problem.ExplanationPath = explanationPath
	}
	return parsedProblem{value: problem, schedules: matter.SubmissionableSchedules}, true, nil
}

func problemSchedule(problem parsedProblem, sections []core.Section, secondaryQualifier bool) (string, error) {
	switch {
	case secondaryQualifier && len(problem.schedules) > 0:
		return "", invalidManifest("problem %q: submissionable_schedules is not allowed in secondary qualifier mode", problem.value.Code)
	case secondaryQualifier:
		return sections[0].Slug, nil
	case len(problem.schedules) == 0:
		return "", invalidManifest("problem %q: submissionable_schedules is required", problem.value.Code)
	case len(problem.schedules) > 1:
		return "", invalidManifest("problem %q: more than one submissionable schedule is not supported yet", problem.value.Code)
	default:
		return problem.schedules[0], nil
	}
}

// redeployRule applies the master repository defaults: no rule means
// effectively unlimited redeploys, and a percentage penalty without a
// percentage costs 10%.
func redeployRule(matter *redeployRuleFrontmatter) (core.RedeployRule, error) {
	if matter == nil {
		threshold, percentage := int32(999), int32(0)
		return core.RedeployRule{Type: core.RedeployPercentage, Threshold: &threshold, Percentage: &percentage}, nil
	}
	rule := core.RedeployRule{Type: core.RedeployRuleType(matter.Type)}
	switch rule.Type {
	case core.RedeployPercentage:
		threshold, percentage := int32(0), int32(10)
		if matter.PenaltyThreshold != nil {
			threshold = *matter.PenaltyThreshold
		}
		if matter.PenaltyPercentage != nil && *matter.PenaltyPercentage != 0 {
			percentage = *matter.PenaltyPercentage
		}
		rule.Threshold, rule.Percentage = &threshold, &percentage
	case core.RedeployUnredeployable, core.RedeployManual:
		// Penalty values have no meaning for these rules and are ignored.
	default:
		return core.RedeployRule{}, fmt.Errorf("unknown type %q", matter.Type)
	}
	return rule, validateRedeployRule(rule)
}

func buildNotice(filePath string, data []byte) (core.Announcement, error) {
	var matter noticeFrontmatter
	body, err := splitFrontmatter(data, &matter)
	if err != nil {
		return core.Announcement{}, invalidManifest("%s: %v", filePath, err)
	}
	slug := strings.TrimSuffix(path.Base(filePath), ".md")
	if !validIdentifier(slug, 255) {
		return core.Announcement{}, invalidManifest("%s: file name is not a valid slug", filePath)
	}
	if strings.TrimSpace(matter.Title) == "" || !validTextLength(matter.Title, 1, 255) {
		return core.Announcement{}, invalidManifest("%s: title is invalid", filePath)
	}
	from, err := time.Parse(time.RFC3339, matter.From)
	if err != nil {
		return core.Announcement{}, invalidManifest("%s: from must be an RFC 3339 date-time", filePath)
	}
	return core.Announcement{
		Slug: slug, Title: matter.Title, Markdown: strings.TrimSpace(body), MarkdownPath: filePath, EffectiveFrom: from.UTC(),
	}, nil
}

func validateSchedule(sections []core.Section) error {
	sorted := append([]core.Section(nil), sections...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Beginning.Before(sorted[j].Beginning) })
	seen := make(map[string]struct{}, len(sorted))
	for index, section := range sorted {
		if !validIdentifier(section.Slug, 255) {
			return invalidManifest("schedule name %q is invalid", section.Slug)
		}
		if _, exists := seen[section.Slug]; exists {
			return invalidManifest("schedule name %q is duplicated", section.Slug)
		}
		seen[section.Slug] = struct{}{}
		if section.Beginning.IsZero() || !section.Beginning.Before(section.Ending) {
			return invalidManifest("schedule %q must have startAt before endAt", section.Slug)
		}
		if index > 0 && section.Beginning.Before(sorted[index-1].Ending) {
			return invalidManifest("schedules %q and %q overlap", sorted[index-1].Slug, section.Slug)
		}
	}
	return nil
}

// splitFrontmatter decodes the leading "---" YAML block into target and
// returns the Markdown that follows it.
func splitFrontmatter(data []byte, target any) (string, error) {
	if !utf8.Valid(data) {
		return "", errors.New("file must be valid UTF-8")
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return "", errors.New("frontmatter is missing")
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	body := ""
	switch {
	case end >= 0:
		body = rest[end+len("\n---\n"):]
	case strings.HasSuffix(rest, "\n---"):
		end = len(rest) - len("\n---")
	default:
		return "", errors.New("frontmatter is not closed")
	}
	if err := decodeStrictYAML([]byte(rest[:end]), target); err != nil {
		return "", fmt.Errorf("frontmatter: %w", err)
	}
	return strings.TrimLeft(body, "\n"), nil
}

func decodeStrictYAML(data []byte, target any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("multiple YAML documents are not allowed")
	}
	return nil
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
