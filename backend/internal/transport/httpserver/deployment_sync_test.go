package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

type missingDeploymentStatusError struct{}

func (missingDeploymentStatusError) Error() string                  { return "deployment status not found" }
func (missingDeploymentStatusError) DeploymentStatusNotFound() bool { return true }

type missingDeploymentStatusGateway struct{}

func (missingDeploymentStatusGateway) Queue(context.Context, service.DeploymentRequest) error {
	return nil
}

func (missingDeploymentStatusGateway) Status(context.Context, int64, string) (core.DeploymentEvent, error) {
	return core.DeploymentEvent{}, missingDeploymentStatusError{}
}

func TestManualSyncNotFoundReleasesQueuedDeployment(t *testing.T) {
	fixture := newContractFixture(t)
	fixture.service.Deployments = missingDeploymentStatusGateway{}
	if _, err := fixture.store.CreateDeployment(t.Context(), core.Deployment{
		RequestID:     "74d8dc35-8a9c-4da3-9fed-d20a41cc7317",
		TeamCode:      2,
		ProblemCode:   "P1",
		ContentCommit: testContentCommit,
		RequestedAt:   testNow,
		LatestStatus:  core.DeploymentQueued,
	}); err != nil {
		t.Fatal(err)
	}
	admin := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "operator"}, service.AdminTTL)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/deployments/2/P1/sync", nil)
	request.Header.Set("Origin", testOrigin)
	request.AddCookie(admin)
	recorder := httptest.NewRecorder()
	fixture.handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", recorder.Code, recorder.Body.String())
	}
	assertProblem(t, recorder, http.StatusNotFound, "resource_not_found")

	team, problem := int64(2), "P1"
	deployments, err := fixture.store.ListDeployments(t.Context(), core.DeploymentFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		t.Fatal(err)
	}
	if len(deployments) != 1 || deployments[0].LatestStatus != core.DeploymentFailed {
		t.Fatalf("deployments = %#v, want one FAILED deployment", deployments)
	}
	if len(deployments[0].Events) != 2 ||
		deployments[0].Events[0].Status != core.DeploymentDeploying ||
		deployments[0].Events[1].Status != core.DeploymentFailed {
		t.Fatalf("recovery events = %#v, want DEPLOYING then FAILED", deployments[0].Events)
	}
}
