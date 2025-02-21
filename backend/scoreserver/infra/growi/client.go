package growi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type Client struct {
	APIToken string
	BaseURL  *url.URL
	Client   *http.Client
}

type (
	GrowiError struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
	growiRespError struct {
		Errors []GrowiError `json:"errors"`
	}
)

func (e GrowiError) Error() string {
	return e.Message
}

func (e GrowiError) Is(target error) bool {
	tgt, ok := target.(GrowiError)
	return ok && (e.Message == "" || e.Message == tgt.Message) && (e.Status == 0 || e.Status == tgt.Status)
}

var _ error = growiRespError{}

func (e growiRespError) Error() string {
	msgs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		msgs = append(msgs, err.Message)
	}
	return strings.Join(msgs, ", ")
}

func (e growiRespError) Unwrap() []error {
	if len(e.Errors) == 0 {
		return nil
	}
	errs := make([]error, len(e.Errors))
	for i, err := range e.Errors {
		errs[i] = err
	}
	return errs
}

const (
	maxResponseSize = 1 << 20
)

func get[V any](ctx context.Context, client *Client, path *url.URL) (*V, error) {
	form := url.Values{}
	form.Add("access_token", client.APIToken)
	reqBody := strings.NewReader(form.Encode())

	path = client.BaseURL.ResolveReference(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path.String(), reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.httpClient().Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}
	respBody := io.LimitReader(resp.Body, maxResponseSize)
	defer func() {
		_, _ = io.Copy(io.Discard, respBody)
		_ = resp.Body.Close()
	}()

	var respErr *domain.Error
	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		respErr = domain.ErrNotFound
	case http.StatusBadRequest:
		respErr = domain.ErrInvalidArgument
	default:
		respErr = domain.ErrInternal
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, errors.WithStack(respErr)
	}
	decoder := json.NewDecoder(respBody)

	if respErr != nil {
		var errs growiRespError
		if err := decoder.Decode(&errs); err != nil {
			return nil, errors.Wrap(err, "failed to decode error response")
		}
		return nil, errors.Join(respErr, errs)
	}

	var v V
	if err := decoder.Decode(&v); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &v, nil
}

func (c *Client) httpClient() *http.Client {
	if c.Client != nil {
		return c.Client
	}
	return http.DefaultClient
}
