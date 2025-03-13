package sstate

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SStateClient struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
	User       string
	Password   string
}

func NewSStateClient(cfg config.SState) (*SStateClient, error) {
	baseURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse base URL")
	}

	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		transport = &http.Transport{}
	}
	transport = transport.Clone()
	if cfg.CA != "" {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM([]byte(cfg.CA)) {
			return nil, errors.New("failed to append CA to cert pool")
		}
		transport.TLSClientConfig.RootCAs = certPool
	}
	if cfg.InsecureSkipVerify {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	return &SStateClient{
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(transport),
		},
		BaseURL:  baseURL,
		User:     cfg.User,
		Password: cfg.Password,
	}, nil
}

var (
	ErrProblemNotFound  = errors.New("problem not found")
	ErrAlreadyDeploying = errors.New("already deploying")
)

type (
	redeployRequest struct {
		TeamID    string `json:"team_id"`
		ProblemID string `json:"problem_id"`
	}
	redeployResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
)

func (c *SStateClient) Redeploy(ctx context.Context, teamCode int64, problemCode string) error {
	ctx, span := tracer.Start(ctx, "SStateClient.Redeploy", trace.WithAttributes(
		attribute.Int64("team_code", teamCode),
		attribute.String("problem_code", problemCode),
	))
	defer span.End()

	reqJSON, err := json.Marshal(redeployRequest{
		TeamID:    fmt.Sprintf("%02d", teamCode),
		ProblemID: problemCode,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal request")
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		c.BaseURL.ResolveReference(&url.URL{Path: "/redeploy"}).String(),
		bytes.NewReader(reqJSON),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.User, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrProblemNotFound
	case http.StatusTooManyRequests:
		return ErrAlreadyDeploying
	}

	var redeployResp redeployResponse
	if err := json.NewDecoder(resp.Body).Decode(&redeployResp); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	return errors.Newf("unexpected status code: %d, message: %s", resp.StatusCode, redeployResp.Message)
}

func (c *SStateClient) GetStatus(ctx context.Context, teamCode int64, problemCode string) (domain.DeploymentStatus, string, error) {
	ctx, span := tracer.Start(ctx, "SStateClient.GetStatus", trace.WithAttributes(
		attribute.Int64("team_code", teamCode),
		attribute.String("problem_code", problemCode),
	))
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.BaseURL.ResolveReference(&url.URL{Path: fmt.Sprintf("/status/%02d/%s", teamCode, problemCode)}).String(), nil,
	)
	if err != nil {
		return domain.DeploymentStatusUnknown, "", errors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(c.User, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return domain.DeploymentStatusUnknown, "", errors.Wrap(err, "failed to send request")
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return domain.DeploymentStatusUnknown, "", ErrProblemNotFound
	}

	var respJSON struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return domain.DeploymentStatusUnknown, "", errors.Wrap(err, "failed to decode response")
	}
	switch respJSON.Status {
	case "Running":
		return domain.DeploymentStatusCompleted, respJSON.Message, nil
	case "Creating", "Queuing":
		return domain.DeploymentStatusCreating, respJSON.Message, nil
	case "Error":
		return domain.DeploymentStatusFailed, respJSON.Message, nil
	default:
		return domain.DeploymentStatusUnknown, respJSON.Message, errors.Newf("unexpected status: %s, message: %s", respJSON.Status, respJSON.Message)
	}
}
