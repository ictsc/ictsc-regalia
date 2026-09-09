package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

type upstreamMachineTokenError struct{}

func (upstreamMachineTokenError) Error() string              { return "identity provider unavailable" }
func (upstreamMachineTokenError) MachineTokenUpstream() bool { return true }

type upstreamMachineTokenContent struct{}

func (upstreamMachineTokenContent) VerifyActionsToken(context.Context, string) (service.ActionsClaims, error) {
	return service.ActionsClaims{}, upstreamMachineTokenError{}
}

func (upstreamMachineTokenContent) Fetch(context.Context, string, string, string) (core.ContentSnapshot, error) {
	panic("Fetch must not be called when authentication fails")
}

func (upstreamMachineTokenContent) IsAncestor(context.Context, string, string, string) (bool, error) {
	panic("IsAncestor must not be called when authentication fails")
}

func TestGitHubActionsIdentityProviderFailureIsBadGateway(t *testing.T) {
	fixture := newContractFixture(t)
	fixture.service.Content = upstreamMachineTokenContent{}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/content/actions/refresh",
		strings.NewReader(`{"commit":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer syntactically-valid-machine-token")
	recorder := httptest.NewRecorder()

	fixture.handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusBadGateway, recorder.Body.String())
	}
	assertProblem(t, recorder, http.StatusBadGateway, "upstream_unavailable")
}
