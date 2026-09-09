package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"gopkg.in/yaml.v3"
)

const (
	maxManifestBytes      = 1 << 20
	maxMarkdownBytes      = 2 << 20
	maxGitHubResponseBody = 4 << 20
)

var (
	ErrInvalidManifest = errors.New("invalid content manifest")
	commitPattern      = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
	identifierPattern  = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	repositoryPart     = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
)

type HTTPError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("GitHub %s returned HTTP %d", e.Operation, e.StatusCode)
}

type contentResponse struct {
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
	Content  string `json:"content"`
	Size     int64  `json:"size"`
	SHA      string `json:"sha"`
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
	manifestBytes, err := c.fetchFile(ctx, owner, name, commit, c.manifestPath, maxManifestBytes)
	if err != nil {
		return core.ContentSnapshot{}, fmt.Errorf("fetch content manifest: %w", err)
	}
	manifest, err := ParseManifest(manifestBytes)
	if err != nil {
		return core.ContentSnapshot{}, err
	}
	files := make(map[string][]byte)
	load := func(relativePath string, maximum int64) ([]byte, error) {
		resolved, err := resolveContentPath(c.contentRoot, relativePath)
		if err != nil {
			return nil, invalidManifest("content path %q: %v", relativePath, err)
		}
		if cached, ok := files[resolved]; ok {
			return append([]byte(nil), cached...), nil
		}
		content, err := c.fetchFile(ctx, owner, name, commit, resolved, maximum)
		if err != nil {
			return nil, err
		}
		files[resolved] = append([]byte(nil), content...)
		return content, nil
	}
	for index := range manifest.Problems {
		problem := &manifest.Problems[index]
		body, err := load(problem.BodyPath, maxMarkdownBytes)
		if err != nil {
			return core.ContentSnapshot{}, wrapContentLoadError(fmt.Sprintf("problem %q body", problem.Code), err)
		}
		if strings.TrimSpace(string(body)) == "" {
			return core.ContentSnapshot{}, invalidManifest("problem %q body is empty", problem.Code)
		}
		problem.Body = string(body)
		if problem.ExplanationPath != "" {
			explanation, err := load(problem.ExplanationPath, maxMarkdownBytes)
			if err != nil {
				return core.ContentSnapshot{}, wrapContentLoadError(fmt.Sprintf("problem %q explanation", problem.Code), err)
			}
			problem.Explanation = string(explanation)
		}
	}
	for index := range manifest.Announcements {
		announcement := &manifest.Announcements[index]
		markdown, err := load(announcement.MarkdownPath, maxMarkdownBytes)
		if err != nil {
			return core.ContentSnapshot{}, wrapContentLoadError(fmt.Sprintf("announcement %q markdown", announcement.Slug), err)
		}
		announcement.Markdown = string(markdown)
	}
	if manifest.RulePath != "" {
		rule, err := load(manifest.RulePath, maxMarkdownBytes)
		if err != nil {
			return core.ContentSnapshot{}, wrapContentLoadError("rule markdown", err)
		}
		manifest.RuleMarkdown = string(rule)
	}
	return core.ContentSnapshot{
		CommitSHA:  commit,
		Repository: repository,
		Ref:        ref,
		Manifest:   manifest,
		FetchedAt:  c.now().UTC(),
	}, nil
}

func (c *Client) IsAncestor(ctx context.Context, repository, ancestor, descendant string) (bool, error) {
	if repository != c.repository {
		return false, errors.New("GitHub repository is outside the configured content source")
	}
	if !commitPattern.MatchString(ancestor) || !commitPattern.MatchString(descendant) {
		return false, errors.New("GitHub compare requires full lowercase object IDs")
	}
	if ancestor == descendant {
		return true, nil
	}
	owner, name, err := splitRepository(repository)
	if err != nil {
		return false, err
	}
	endpoint := fmt.Sprintf("%s/repos/%s/%s/compare/%s...%s", c.apiBaseURL,
		url.PathEscape(owner), url.PathEscape(name), url.PathEscape(ancestor), url.PathEscape(descendant))
	var response struct {
		Status string `json:"status"`
	}
	if err := c.getGitHubJSON(ctx, endpoint, "compare commits", &response); err != nil {
		return false, err
	}
	switch response.Status {
	case "ahead", "identical":
		return true, nil
	case "behind", "diverged":
		return false, nil
	default:
		return false, errors.New("GitHub compare response contains an unknown status")
	}
}

func ParseManifest(data []byte) (core.Manifest, error) {
	if len(data) == 0 || len(data) > maxManifestBytes {
		return core.Manifest{}, invalidManifest("manifest size is invalid")
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var manifest core.Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return core.Manifest{}, invalidManifest("decode YAML: %v", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return core.Manifest{}, invalidManifest("multiple YAML documents are not allowed")
		}
		return core.Manifest{}, invalidManifest("decode trailing YAML: %v", err)
	}
	if err := validateManifest(manifest); err != nil {
		return core.Manifest{}, err
	}
	return manifest, nil
}

func validateManifest(manifest core.Manifest) error {
	if manifest.Version != 1 {
		return invalidManifest("version must be 1")
	}
	if manifest.RuleMarkdown != "" {
		return invalidManifest("inline rule_markdown is not allowed; use rule_path")
	}
	if manifest.RulePath != "" && !safeContentRelativePath(manifest.RulePath) {
		return invalidManifest("rule_path is unsafe")
	}
	sections := append([]core.Section(nil), manifest.Sections...)
	sort.Slice(sections, func(i, j int) bool {
		if sections[i].Beginning.Equal(sections[j].Beginning) {
			return sections[i].Slug < sections[j].Slug
		}
		return sections[i].Beginning.Before(sections[j].Beginning)
	})
	sectionBySlug := make(map[string]core.Section, len(sections))
	for index, section := range sections {
		if !validIdentifier(section.Slug, 255) {
			return invalidManifest("section slug %q is invalid", section.Slug)
		}
		if section.Beginning.IsZero() || section.Ending.IsZero() || !section.Beginning.Before(section.Ending) {
			return invalidManifest("section %q must have a non-empty [beginning, ending) interval", section.Slug)
		}
		if _, exists := sectionBySlug[section.Slug]; exists {
			return invalidManifest("section slug %q is duplicated", section.Slug)
		}
		if index > 0 && section.Beginning.Before(sections[index-1].Ending) {
			return invalidManifest("sections %q and %q overlap", sections[index-1].Slug, section.Slug)
		}
		sectionBySlug[section.Slug] = section
	}
	problemByCode := make(map[string]core.Problem, len(manifest.Problems))
	for _, problem := range manifest.Problems {
		if !validIdentifier(problem.Code, 8) {
			return invalidManifest("problem code %q is invalid", problem.Code)
		}
		if _, exists := problemByCode[problem.Code]; exists {
			return invalidManifest("problem code %q is duplicated", problem.Code)
		}
		if strings.TrimSpace(problem.Title) == "" || !validTextLength(problem.Title, 1, 255) ||
			!validTextLength(problem.Category, 0, 255) {
			return invalidManifest("problem %q title or category length is invalid", problem.Code)
		}
		if problem.MaxScore <= 0 || problem.Type != "DESCRIPTIVE" {
			return invalidManifest("problem %q max_score or type is invalid", problem.Code)
		}
		if _, exists := sectionBySlug[problem.SectionSlug]; !exists {
			return invalidManifest("problem %q references unknown section %q", problem.Code, problem.SectionSlug)
		}
		if problem.Body != "" || problem.Explanation != "" {
			return invalidManifest("problem %q may not contain inline markdown", problem.Code)
		}
		if !safeContentRelativePath(problem.BodyPath) {
			return invalidManifest("problem %q body_path is unsafe", problem.Code)
		}
		if problem.ExplanationPath != "" && !safeContentRelativePath(problem.ExplanationPath) {
			return invalidManifest("problem %q explanation_path is unsafe", problem.Code)
		}
		if err := validateRedeployRule(problem.Redeploy); err != nil {
			return invalidManifest("problem %q redeploy_rule: %v", problem.Code, err)
		}
		problemByCode[problem.Code] = problem
	}
	problemSectionCount := make(map[string]int, len(problemByCode))
	for _, section := range manifest.Sections {
		seenInSection := make(map[string]struct{}, len(section.ProblemIDs))
		for _, code := range section.ProblemIDs {
			problem, exists := problemByCode[code]
			if !exists {
				return invalidManifest("section %q references unknown problem %q", section.Slug, code)
			}
			if _, duplicate := seenInSection[code]; duplicate {
				return invalidManifest("section %q repeats problem %q", section.Slug, code)
			}
			if problem.SectionSlug != section.Slug {
				return invalidManifest("problem %q section_slug does not match section %q", code, section.Slug)
			}
			seenInSection[code] = struct{}{}
			problemSectionCount[code]++
		}
	}
	for code := range problemByCode {
		if problemSectionCount[code] != 1 {
			return invalidManifest("problem %q must occur in exactly one section problems list", code)
		}
	}
	announcementSlugs := make(map[string]struct{}, len(manifest.Announcements))
	for _, announcement := range manifest.Announcements {
		if !validIdentifier(announcement.Slug, 255) || strings.TrimSpace(announcement.Title) == "" ||
			!validTextLength(announcement.Title, 1, 255) {
			return invalidManifest("announcement %q slug or title is invalid", announcement.Slug)
		}
		if _, exists := announcementSlugs[announcement.Slug]; exists {
			return invalidManifest("announcement slug %q is duplicated", announcement.Slug)
		}
		if announcement.EffectiveFrom.IsZero() {
			return invalidManifest("announcement %q effective_from is required", announcement.Slug)
		}
		if announcement.Markdown != "" {
			return invalidManifest("announcement %q may not contain inline markdown", announcement.Slug)
		}
		if !safeContentRelativePath(announcement.MarkdownPath) {
			return invalidManifest("announcement %q markdown_path is unsafe", announcement.Slug)
		}
		announcementSlugs[announcement.Slug] = struct{}{}
	}
	return nil
}

func validateRedeployRule(rule core.RedeployRule) error {
	switch rule.Type {
	case core.RedeployPercentage:
		if rule.Threshold == nil || rule.Percentage == nil {
			return errors.New("percentage penalty requires penalty_threshold and penalty_percentage")
		}
		if *rule.Threshold < 0 || *rule.Percentage < 0 || *rule.Percentage > 99 {
			return errors.New("percentage penalty values are outside their allowed ranges")
		}
	case core.RedeployManual, core.RedeployUnredeployable:
		if rule.Threshold != nil || rule.Percentage != nil {
			return errors.New("non-percentage rule requires null penalty values")
		}
	default:
		return errors.New("unknown redeploy rule type")
	}
	return nil
}

func (c *Client) verifyCommit(ctx context.Context, owner, repository, commit string) error {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/git/commits/%s", c.apiBaseURL,
		url.PathEscape(owner), url.PathEscape(repository), url.PathEscape(commit))
	var response struct {
		SHA string `json:"sha"`
	}
	if err := c.getGitHubJSON(ctx, endpoint, "get commit", &response); err != nil {
		return err
	}
	if response.SHA != commit {
		return errors.New("GitHub commit response does not match the requested object ID")
	}
	return nil
}

func (c *Client) fetchFile(ctx context.Context, owner, repository, commit, filePath string, maximum int64) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", c.apiBaseURL,
		url.PathEscape(owner), url.PathEscape(repository), escapeRepositoryPath(filePath), url.QueryEscape(commit))
	var response contentResponse
	if err := c.getGitHubJSON(ctx, endpoint, "get repository content", &response); err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return nil, invalidContentError{cause: err}
		}
		return nil, err
	}
	if response.Type != "file" || response.Encoding != "base64" || response.Size < 0 || response.Size > maximum {
		return nil, invalidContentError{cause: errors.New("GitHub content response does not describe an allowed file")}
	}
	decoded, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return nil, invalidContentError{cause: errors.New("GitHub content response has invalid base64")}
	}
	if int64(len(decoded)) != response.Size || int64(len(decoded)) > maximum {
		return nil, invalidContentError{cause: errors.New("GitHub content response size does not match")}
	}
	if !utf8.Valid(decoded) {
		return nil, invalidContentError{cause: errors.New("GitHub content file must be valid UTF-8")}
	}
	return decoded, nil
}

func (c *Client) getGitHubJSON(ctx context.Context, endpoint, operation string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create GitHub %s request: %w", operation, err)
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("User-Agent", "ictsc-regalia")
	if c.apiToken != "" {
		request.Header.Set("Authorization", "Bearer "+c.apiToken)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("GitHub %s: %w", operation, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return &HTTPError{Operation: operation, StatusCode: response.StatusCode, Body: strings.TrimSpace(string(body))}
	}
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxGitHubResponseBody+1))
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode GitHub %s response: %w", operation, err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return fmt.Errorf("decode GitHub %s response: %w", operation, err)
	}
	return nil
}

func splitRepository(repository string) (string, string, error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || !repositoryPart.MatchString(parts[0]) || !repositoryPart.MatchString(parts[1]) ||
		parts[0] == "." || parts[0] == ".." || parts[1] == "." || parts[1] == ".." {
		return "", "", errors.New("GitHub repository must have the form owner/name")
	}
	return parts[0], parts[1], nil
}

func safeRepositoryPath(value string) bool {
	if value == "" || strings.Contains(value, "\\") || strings.HasPrefix(value, "/") {
		return false
	}
	cleaned := path.Clean(value)
	return cleaned == value && cleaned != "." && !strings.HasPrefix(cleaned, "../")
}

func safeContentRelativePath(value string) bool {
	return safeRepositoryPath(value)
}

func resolveContentPath(root, relative string) (string, error) {
	if !safeContentRelativePath(relative) {
		return "", errors.New("path must be clean and relative to the content root")
	}
	if root == "." {
		return relative, nil
	}
	resolved := path.Join(root, relative)
	if !strings.HasPrefix(resolved, root+"/") {
		return "", errors.New("path escapes the content root")
	}
	return resolved, nil
}

func repositoryDir(filePath string) string {
	directory := path.Dir(filePath)
	if directory == "" {
		return "."
	}
	return directory
}

func escapeRepositoryPath(value string) string {
	parts := strings.Split(value, "/")
	for index := range parts {
		parts[index] = url.PathEscape(parts[index])
	}
	return strings.Join(parts, "/")
}

func validIdentifier(value string, maximum int) bool {
	return len(value) >= 1 && len(value) <= maximum && identifierPattern.MatchString(value)
}

func validTextLength(value string, minimum, maximum int) bool {
	if !utf8.ValidString(value) {
		return false
	}
	length := utf8.RuneCountInString(value)
	return length >= minimum && length <= maximum
}

func invalidManifest(format string, values ...any) error {
	return invalidContentError{cause: fmt.Errorf("%w: %s", ErrInvalidManifest, fmt.Sprintf(format, values...))}
}
