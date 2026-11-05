package rstate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	defaultTimeout = 15 * time.Second
	// rstate's admin routes use gin.BasicAuth with this fixed account.
	defaultAdminUser = "ictsc"
	maxResponseBytes = 4 << 20
)

// rstate rejects problem IDs shorter than three characters.
var problemCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{3,8}$`)

var (
	_ service.DeploymentGateway = (*Client)(nil)
	_ service.DeploymentProber  = (*Client)(nil)
)

type Config struct {
	BaseURL string
	// AdminUser defaults to rstate's fixed "ictsc" account; set it when a
	// proxy in front of rstate checks other credentials.
	AdminUser         string
	AdminPassword     string
	HTTPClient        *http.Client
	AllowInsecureHTTP bool
}

type Client struct {
	baseURL       *url.URL
	adminUser     string
	adminPassword string
	httpClient    *http.Client
}

// Job is one entry of rstate's GET /admin/list. TeamID has the "team" prefix
// stripped by rstate, and both times are Unix nanoseconds.
type Job struct {
	ID          string `json:"id"`
	State       string `json:"state"`
	CreatedTime int64  `json:"created_time"`
	EndTime     int64  `json:"end_time"`
	TeamID      string `json:"team_id"`
	ProbID      string `json:"prob_id"`
}

// rstate's GetJobList reports states as these display strings.
const (
	stateWait    = "開始待ち"
	stateRunning = "実行中"
	stateSuccess = "終了"
	stateError   = "エラー"
	stateLimit   = "Limit"
)

// postJobBadRequest is the body rstate's postJob sends for an invalid team or
// problem ID; its other 400s carry the team's status page URL instead.
const postJobBadRequest = "BadRequest!"

type HTTPError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("rstate %s returned HTTP %d", e.Operation, e.StatusCode)
	}
	return fmt.Sprintf("rstate %s returned HTTP %d: %s", e.Operation, e.StatusCode, e.Body)
}

// TargetError means rstate cannot take the team or problem at all, so the
// request is rejected before it is sent and is never worth retrying.
type TargetError struct {
	Reason string
}

func (e *TargetError) Error() string { return e.Reason }

// NotFoundError means rstate has never run a job for the team's problem.
type NotFoundError struct {
	TeamCode    int64
	ProblemCode string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("rstate has no job for %s/%s", teamID(e.TeamCode), e.ProblemCode)
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.AdminPassword) == "" {
		return nil, errors.New("rstate admin password is required")
	}
	parsed, err := url.Parse(config.BaseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil {
		return nil, errors.New("rstate base URL must be absolute and may not include user information")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("rstate base URL may not include a query or fragment")
	}
	if parsed.Scheme != "https" && !(config.AllowInsecureHTTP && parsed.Scheme == "http") {
		return nil, errors.New("rstate base URL must use HTTPS")
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: defaultTimeout}
	}
	if config.AdminUser == "" {
		config.AdminUser = defaultAdminUser
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return &Client{baseURL: parsed, adminUser: config.AdminUser, adminPassword: config.AdminPassword, httpClient: config.HTTPClient}, nil
}

// Queue uses POST /admin/postJob behind Basic auth, as the ICTSC Discord bot
// did. rstate refuses contestant teams there until its preparatory-phase
// dates (pkg/utils/util.go) are set for the contest.
func (c *Client) Queue(ctx context.Context, request service.DeploymentRequest) (string, error) {
	if err := validateTarget(request.TeamCode, request.ProblemCode); err != nil {
		return "", err
	}
	form := url.Values{"team_id": {teamID(request.TeamCode)}, "prob_id": {probID(request.ProblemCode)}}
	response, err := c.do(ctx, http.MethodPost, c.endpoint("/admin/postJob"), strings.NewReader(form.Encode()), "queue redeployment")
	if err != nil {
		return "", err
	}
	defer drainAndClose(response.Body)
	switch response.StatusCode {
	case http.StatusOK:
		var job struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(io.LimitReader(response.Body, maxResponseBytes)).Decode(&job); err != nil || job.ID == "" {
			// rstate has taken the job; polling adopts it as the team's
			// running job without the ID.
			return "", nil
		}
		return job.ID, nil
	case http.StatusBadRequest:
		// The body is either "BadRequest!" or the team's status page URL,
		// which carries a token and must not reach deployment messages.
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		if strings.TrimSpace(string(body)) == postJobBadRequest {
			return "", &HTTPError{Operation: "queue redeployment", StatusCode: response.StatusCode, Body: "rstate rejected the team or problem ID"}
		}
		return c.resolveRefusedJob(ctx, request)
	default:
		return "", readHTTPError("queue redeployment", response)
	}
}

// resolveRefusedJob settles postJob's ambiguous 400, sent both when the
// team's problem already has a waiting or running job and in rstate's
// preparatory phase. A retry after a lost response lands here too, so an
// active job means the request was taken and returns that job's ID.
func (c *Client) resolveRefusedJob(ctx context.Context, request service.DeploymentRequest) (string, error) {
	jobs, err := c.ListJobs(ctx)
	if err != nil {
		return "", fmt.Errorf("check rstate jobs after a refused redeployment: %w", err)
	}
	for _, job := range jobs {
		if job.TeamCode == request.TeamCode && job.ProblemCode == probID(request.ProblemCode) &&
			(job.Status == core.DeploymentQueued || job.Status == core.DeploymentDeploying) {
			return job.ID, nil
		}
	}
	return "", &HTTPError{Operation: "queue redeployment", StatusCode: http.StatusBadRequest, Body: "rstate is in its preparatory phase"}
}

// ListJobs returns every job rstate reports, skipping entries whose team or
// state Regalia cannot read.
func (c *Client) ListJobs(ctx context.Context) ([]service.DeploymentJob, error) {
	raw, err := c.listJobs(ctx)
	if err != nil {
		return nil, err
	}
	jobs := make([]service.DeploymentJob, 0, len(raw))
	for _, entry := range raw {
		teamCode, err := strconv.ParseInt(entry.TeamID, 10, 64)
		if err != nil {
			continue
		}
		status, err := mapState(entry.State)
		if err != nil {
			continue
		}
		job := service.DeploymentJob{ID: entry.ID, TeamCode: teamCode, ProblemCode: entry.ProbID, Status: status}
		if status == core.DeploymentFailed {
			message := fmt.Sprintf("rstate job %s ended with state %s", entry.ID, entry.State)
			job.Message = &message
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Status reports the newest job for the team's problem. The operator sync
// does not use it for rstate, which goes through the service's polling path.
func (c *Client) Status(ctx context.Context, teamCode int64, problemCode string) (core.DeploymentEvent, error) {
	if err := validateTarget(teamCode, problemCode); err != nil {
		return core.DeploymentEvent{}, err
	}
	raw, err := c.listJobs(ctx)
	if err != nil {
		return core.DeploymentEvent{}, err
	}
	team, problem := fmt.Sprintf("%02d", teamCode), probID(problemCode)
	var latest *Job
	for i := range raw {
		if raw[i].TeamID == team && raw[i].ProbID == problem && (latest == nil || raw[i].CreatedTime > latest.CreatedTime) {
			latest = &raw[i]
		}
	}
	if latest == nil {
		return core.DeploymentEvent{}, &NotFoundError{TeamCode: teamCode, ProblemCode: problemCode}
	}
	status, err := mapState(latest.State)
	if err != nil {
		return core.DeploymentEvent{}, err
	}
	if status == core.DeploymentQueued {
		status = core.DeploymentDeploying
	}
	return core.DeploymentEvent{Status: status}, nil
}

func (c *Client) listJobs(ctx context.Context) ([]Job, error) {
	response, err := c.do(ctx, http.MethodGet, c.endpoint("/admin/list"), nil, "list jobs")
	if err != nil {
		return nil, err
	}
	defer drainAndClose(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, readHTTPError("list jobs", response)
	}
	var jobs []Job
	if err := json.NewDecoder(io.LimitReader(response.Body, maxResponseBytes)).Decode(&jobs); err != nil {
		return nil, fmt.Errorf("decode rstate job list: %w", err)
	}
	return jobs, nil
}

func (c *Client) do(ctx context.Context, method, endpoint string, form io.Reader, operation string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, endpoint, form)
	if err != nil {
		return nil, fmt.Errorf("create rstate %s request: %w", operation, err)
	}
	request.SetBasicAuth(c.adminUser, c.adminPassword)
	request.Header.Set("Accept", "application/json")
	if form != nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("rstate %s: %w", operation, err)
	}
	return response, nil
}

func (c *Client) endpoint(endpointPath string) string {
	copy := *c.baseURL
	copy.Path = c.baseURL.Path + endpointPath
	return copy.String()
}

func teamID(teamCode int64) string {
	return fmt.Sprintf("team%02d", teamCode)
}

// probID lowercases the problem code: rstate's Terraform modules are named
// in lower case, and both the Discord bot and rstate's SState-compatible API
// send codes that way.
func probID(problemCode string) string {
	return strings.ToLower(problemCode)
}

func validateTarget(teamCode int64, problemCode string) error {
	// rstate only needs the "team" prefix; 99 keeps teamID's two digits in
	// line with rstate's team01 to team99 workspaces.
	if teamCode < 1 || teamCode > 99 {
		return &TargetError{Reason: "rstate team code must be between 1 and 99"}
	}
	if !problemCodePattern.MatchString(problemCode) {
		return &TargetError{Reason: "rstate problem code must be 3 to 8 letters, digits, '_' or '-'"}
	}
	return nil
}

func mapState(state string) (core.DeploymentStatus, error) {
	switch state {
	case stateWait:
		return core.DeploymentQueued, nil
	case stateRunning:
		return core.DeploymentDeploying, nil
	case stateSuccess:
		return core.DeploymentCompleted, nil
	case stateError, stateLimit:
		return core.DeploymentFailed, nil
	default:
		return "", fmt.Errorf("rstate reported an unsupported job state %q", state)
	}
}

func readHTTPError(operation string, response *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
	return &HTTPError{Operation: operation, StatusCode: response.StatusCode, Body: strings.TrimSpace(string(body))}
}

func drainAndClose(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, io.LimitReader(body, maxResponseBytes))
	_ = body.Close()
}
