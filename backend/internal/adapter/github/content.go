package github

import (
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
	"strings"
	"unicode/utf8"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

const (
	maxScheduleBytes      = 1 << 20
	maxMarkdownBytes      = 2 << 20
	maxImageBytes         = 10 << 20
	maxTotalImageBytes    = 100 << 20
	maxGitHubResponseBody = 16 << 20
)

var (
	ErrInvalidManifest = errors.New("invalid content repository")
	commitPattern      = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
	identifierPattern  = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	repositoryPart     = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	discordIDPattern   = regexp.MustCompile(`^[0-9]+$`)
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
