package sstate

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	defaultTimeout   = 15 * time.Second
	maxResponseBytes = 64 << 10
)

var (
	problemCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,8}$`)
	commitPattern      = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
	uuidPattern        = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

var _ service.DeploymentGateway = (*Client)(nil)

type Config struct {
	BaseURL           string
	BearerToken       string
	HTTPClient        *http.Client
	AllowInsecureHTTP bool
}

type Client struct {
	baseURL     *url.URL
	bearerToken string
	httpClient  *http.Client
}

type RedeployRequest struct {
	RequestID     string `json:"request_id"`
	TeamCode      int64  `json:"team_code"`
	ProblemCode   string `json:"problem_code"`
	Revision      int32  `json:"revision"`
	ContentCommit string `json:"content_commit"`
	CallbackURL   string `json:"callback_url"`
}

type StatusResponse struct {
	EventID    string                `json:"event_id"`
	OccurredAt time.Time             `json:"occurred_at"`
	Status     core.DeploymentStatus `json:"status"`
	Message    *string               `json:"message"`
}

type HTTPError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("SState %s returned HTTP %d", e.Operation, e.StatusCode)
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.BearerToken) == "" {
		return nil, errors.New("SState bearer token is required")
	}
	parsed, err := url.Parse(config.BaseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil {
		return nil, errors.New("SState base URL must be absolute and may not include user information")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("SState base URL may not include a query or fragment")
	}
	if parsed.Scheme != "https" && !(config.AllowInsecureHTTP && parsed.Scheme == "http") {
		return nil, errors.New("SState base URL must use HTTPS")
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: defaultTimeout}
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return &Client{baseURL: parsed, bearerToken: config.BearerToken, httpClient: config.HTTPClient}, nil
}

func (c *Client) Queue(ctx context.Context, request service.DeploymentRequest) error {
	payload := RedeployRequest{
		RequestID: request.RequestID, TeamCode: request.TeamCode, ProblemCode: request.ProblemCode,
		Revision: request.Revision, ContentCommit: request.ContentCommit, CallbackURL: request.CallbackURL,
	}
	if err := validateRedeployRequest(payload); err != nil {
		return err
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode SState redeploy request: %w", err)
	}
	response, err := c.do(ctx, http.MethodPost, c.endpoint("/redeploy"), bytes.NewReader(encoded), "queue redeployment")
	if err != nil {
		return err
	}
	defer drainAndClose(response.Body)
	if response.StatusCode != http.StatusAccepted {
		return readHTTPError("queue redeployment", response)
	}
	return nil
}

func (c *Client) Status(ctx context.Context, teamCode int64, problemCode string) (core.DeploymentEvent, error) {
	if teamCode < 2 || teamCode > 99 {
		return core.DeploymentEvent{}, errors.New("SState team code must be between 2 and 99")
	}
	if !problemCodePattern.MatchString(problemCode) {
		return core.DeploymentEvent{}, errors.New("SState problem code is invalid")
	}
	endpoint := c.endpoint(fmt.Sprintf("/status/%02d/%s", teamCode, url.PathEscape(problemCode)))
	response, err := c.do(ctx, http.MethodGet, endpoint, nil, "get redeployment status")
	if err != nil {
		return core.DeploymentEvent{}, err
	}
	defer drainAndClose(response.Body)
	if response.StatusCode != http.StatusOK {
		return core.DeploymentEvent{}, readHTTPError("get redeployment status", response)
	}
	var status StatusResponse
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxResponseBytes+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&status); err != nil {
		return core.DeploymentEvent{}, fmt.Errorf("decode SState status response: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return core.DeploymentEvent{}, fmt.Errorf("decode SState status response: %w", err)
	}
	if err := validateStatus(status); err != nil {
		return core.DeploymentEvent{}, err
	}
	return core.DeploymentEvent{
		EventID: status.EventID, OccurredAt: status.OccurredAt.UTC(), Status: status.Status, Message: status.Message,
	}, nil
}

func (c *Client) do(ctx context.Context, method, endpoint string, body io.Reader, operation string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create SState %s request: %w", operation, err)
	}
	request.Header.Set("Authorization", "Bearer "+c.bearerToken)
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("SState %s: %w", operation, err)
	}
	return response, nil
}

func (c *Client) endpoint(endpointPath string) string {
	copy := *c.baseURL
	copy.Path = strings.TrimRight(c.baseURL.Path, "/") + endpointPath
	return copy.String()
}

func validateRedeployRequest(request RedeployRequest) error {
	if !uuidPattern.MatchString(request.RequestID) {
		return errors.New("SState request_id must be a lowercase UUID")
	}
	if request.TeamCode < 2 || request.TeamCode > 99 {
		return errors.New("SState team_code must be between 2 and 99")
	}
	if !problemCodePattern.MatchString(request.ProblemCode) {
		return errors.New("SState problem_code is invalid")
	}
	if request.Revision < 1 {
		return errors.New("SState revision must be positive")
	}
	if !commitPattern.MatchString(request.ContentCommit) {
		return errors.New("SState content_commit must be a full lowercase object ID")
	}
	callback, err := url.Parse(request.CallbackURL)
	if err != nil || callback.Scheme == "" || callback.Host == "" || callback.User != nil {
		return errors.New("SState callback_url must be an absolute HTTP(S) URL without user information")
	}
	if callback.Scheme != "https" && callback.Scheme != "http" {
		return errors.New("SState callback_url must use HTTP or HTTPS")
	}
	if callback.Fragment != "" {
		return errors.New("SState callback_url may not include a fragment")
	}
	return nil
}

func validateStatus(status StatusResponse) error {
	if !uuidPattern.MatchString(status.EventID) {
		return errors.New("SState status event_id must be a lowercase UUID")
	}
	if status.OccurredAt.IsZero() {
		return errors.New("SState status occurred_at is required")
	}
	switch status.Status {
	case core.DeploymentDeploying, core.DeploymentCompleted, core.DeploymentFailed:
	default:
		return errors.New("SState status contains an unsupported deployment state")
	}
	if status.Message != nil && utf8.RuneCountInString(*status.Message) > 2000 {
		return errors.New("SState status message exceeds 2000 characters")
	}
	return nil
}

func readHTTPError(operation string, response *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
	return &HTTPError{Operation: operation, StatusCode: response.StatusCode, Body: strings.TrimSpace(string(body))}
}

func requireJSONEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}

func drainAndClose(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, io.LimitReader(body, maxResponseBytes))
	_ = body.Close()
}

func EqualBearerToken(left, right string) bool {
	return len(left) == len(right) && subtle.ConstantTimeCompare([]byte(left), []byte(right)) == 1
}
