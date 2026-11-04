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
	"strings"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	defaultTimeout = 15 * time.Second
	// rstate's admin routes use gin.BasicAuth with this fixed account.
	adminUser        = "ictsc"
	maxResponseBytes = 4 << 20
)

// rstate rejects problem IDs shorter than three characters.
var problemCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{3,8}$`)

var (
	_ service.DeploymentGateway = (*Client)(nil)
	_ service.DeploymentProber  = (*Client)(nil)
)

type Config struct {
	BaseURL           string
	AdminPassword     string
	HTTPClient        *http.Client
	AllowInsecureHTTP bool
}

type Client struct {
	baseURL       *url.URL
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
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return &Client{baseURL: parsed, adminPassword: config.AdminPassword, httpClient: config.HTTPClient}, nil
}

// Queue uses POST /backend/:teamid/:probcode rather than /admin/postJob:
// postJob refuses every contestant team because rstate's preparatory-phase
// check compares against fixed 2022 timestamps.
func (c *Client) Queue(ctx context.Context, request service.DeploymentRequest) error {
	if err := validateTarget(request.TeamCode, request.ProblemCode); err != nil {
		return err
	}
	endpoint := c.endpoint("/backend/" + teamID(request.TeamCode) + "/" + url.PathEscape(request.ProblemCode))
	response, err := c.do(ctx, http.MethodPost, endpoint, "queue redeployment")
	if err != nil {
		return err
	}
	defer drainAndClose(response.Body)
	if response.StatusCode != http.StatusOK {
		return readHTTPError("queue redeployment", response)
	}
	return nil
}

// Observe returns rstate's newest job for the team's problem.
func (c *Client) Observe(ctx context.Context, teamCode int64, problemCode string) (service.DeploymentObservation, error) {
	if err := validateTarget(teamCode, problemCode); err != nil {
		return service.DeploymentObservation{}, err
	}
	jobs, err := c.ListJobs(ctx)
	if err != nil {
		return service.DeploymentObservation{}, err
	}
	team := fmt.Sprintf("%02d", teamCode)
	var latest *Job
	for i := range jobs {
		job := &jobs[i]
		if job.TeamID != team || job.ProbID != problemCode {
			continue
		}
		if latest == nil || job.CreatedTime > latest.CreatedTime {
			latest = job
		}
	}
	if latest == nil {
		return service.DeploymentObservation{}, &NotFoundError{TeamCode: teamCode, ProblemCode: problemCode}
	}
	status, err := mapState(latest.State)
	if err != nil {
		return service.DeploymentObservation{}, err
	}
	observation := service.DeploymentObservation{
		JobID: latest.ID, CreatedAt: time.Unix(0, latest.CreatedTime).UTC(), Status: status,
	}
	if status == core.DeploymentFailed {
		message := fmt.Sprintf("rstate job %s ended with state %s", latest.ID, latest.State)
		observation.Message = &message
	}
	return observation, nil
}

// Status serves the operator sync. rstate keeps no event IDs, so the caller
// assigns the event ID and time; a waiting job is reported as DEPLOYING
// because QUEUED cannot be sent as an event.
func (c *Client) Status(ctx context.Context, teamCode int64, problemCode string) (core.DeploymentEvent, error) {
	observation, err := c.Observe(ctx, teamCode, problemCode)
	if err != nil {
		return core.DeploymentEvent{}, err
	}
	status := observation.Status
	if status == core.DeploymentQueued {
		status = core.DeploymentDeploying
	}
	return core.DeploymentEvent{Status: status, Message: observation.Message}, nil
}

func (c *Client) ListJobs(ctx context.Context) ([]Job, error) {
	response, err := c.do(ctx, http.MethodGet, c.endpoint("/admin/list"), "list jobs")
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

func (c *Client) do(ctx context.Context, method, endpoint, operation string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create rstate %s request: %w", operation, err)
	}
	request.SetBasicAuth(adminUser, c.adminPassword)
	request.Header.Set("Accept", "application/json")
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

func validateTarget(teamCode int64, problemCode string) error {
	if teamCode < 2 || teamCode > 99 {
		return errors.New("rstate team code must be between 2 and 99")
	}
	if !problemCodePattern.MatchString(problemCode) {
		return errors.New("rstate problem code must be 3 to 8 letters, digits, '_' or '-'")
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
