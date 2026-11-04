package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

type fakeProber struct {
	jobs  []service.DeploymentJob
	err   error
	calls int
}

func (*fakeProber) Queue(context.Context, service.DeploymentRequest) (string, error) { return "", nil }

func (*fakeProber) Status(context.Context, int64, string) (core.DeploymentEvent, error) {
	return core.DeploymentEvent{}, nil
}

func (p *fakeProber) ListJobs(context.Context) ([]service.DeploymentJob, error) {
	p.calls++
	return p.jobs, p.err
}

func newPollingService(t *testing.T, now time.Time, prober service.DeploymentGateway) (*service.Service, *memory.CompetitionStore, core.Deployment) {
	t.Helper()
	store := memory.NewCompetitionStore()
	svc := service.New(store, memory.NewSessionStore(), service.Config{})
	svc.Now = func() time.Time { return now }
	svc.Deployments = prober
	deployment := createQueuedDeployment(t, store, "74d8dc35-8a9c-4da3-9fed-d20a41cc7317", 2, "A01", now.Add(-time.Minute))
	return svc, store, deployment
}

func createQueuedDeployment(t *testing.T, store *memory.CompetitionStore, requestID string, team int64, problem string, requestedAt time.Time) core.Deployment {
	t.Helper()
	deployment, err := store.CreateDeployment(t.Context(), core.Deployment{
		RequestID: requestID, TeamCode: team, ProblemCode: problem,
		ContentCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", RequestedAt: requestedAt, LatestStatus: core.DeploymentQueued,
		Events: []core.DeploymentEvent{{EventID: requestID[:8] + "-c9fc-4c58-b46e-17b176bba691", OccurredAt: requestedAt, Status: core.DeploymentQueued}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return deployment
}

func deploymentByID(t *testing.T, store *memory.CompetitionStore, requestID string) core.Deployment {
	t.Helper()
	deployments, err := store.ListDeployments(t.Context(), core.DeploymentFilter{})
	if err != nil {
		t.Fatal(err)
	}
	for _, deployment := range deployments {
		if deployment.RequestID == requestID {
			return deployment
		}
	}
	t.Fatalf("deployment %s not found", requestID)
	return core.Deployment{}
}

func job(id string, status core.DeploymentStatus) service.DeploymentJob {
	return service.DeploymentJob{ID: id, TeamCode: 2, ProblemCode: "a01", Status: status}
}

func TestPollDeploymentsFollowsTheRecordedJobOnly(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	// The earlier job finished moments ago; only "new" belongs to the request.
	prober := &fakeProber{jobs: []service.DeploymentJob{job("old", core.DeploymentFailed), job("new", core.DeploymentQueued)}}
	svc, store, deployment := newPollingService(t, now, prober)
	if err := store.SetDeploymentJobID(t.Context(), deployment.RequestID, "new"); err != nil {
		t.Fatal(err)
	}

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := deploymentByID(t, store, deployment.RequestID); got.LatestStatus != core.DeploymentQueued || len(got.Events) != 1 {
		t.Fatalf("the earlier job's failure reached the deployment: %#v", got)
	}

	prober.jobs = []service.DeploymentJob{job("old", core.DeploymentFailed), job("new", core.DeploymentCompleted)}
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	got := deploymentByID(t, store, deployment.RequestID)
	if got.LatestStatus != core.DeploymentCompleted || len(got.Events) != 3 {
		t.Fatalf("deployment = %#v, want QUEUED -> DEPLOYING -> COMPLETED", got)
	}
	if !got.Events[2].OccurredAt.After(got.Events[1].OccurredAt) {
		t.Fatalf("events are not strictly ordered: %#v", got.Events)
	}
}

func TestPollDeploymentsAdoptsOnlyAnActiveJobWithoutID(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{jobs: []service.DeploymentJob{job("old", core.DeploymentCompleted)}}
	svc, store, deployment := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := deploymentByID(t, store, deployment.RequestID); got.LatestStatus != core.DeploymentQueued || got.GatewayJobID != "" {
		t.Fatalf("a finished job was adopted: %#v", got)
	}

	prober.jobs = append(prober.jobs, job("running", core.DeploymentDeploying))
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := deploymentByID(t, store, deployment.RequestID); got.LatestStatus != core.DeploymentDeploying || got.GatewayJobID != "running" {
		t.Fatalf("deployment = %#v, want the running job adopted", got)
	}
}

func TestPollDeploymentsKeepsWaitingAndRepeatedStatesQuiet(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{jobs: []service.DeploymentJob{job("job", core.DeploymentQueued)}}
	svc, store, deployment := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := deploymentByID(t, store, deployment.RequestID); got.LatestStatus != core.DeploymentQueued || len(got.Events) != 1 {
		t.Fatalf("waiting job changed the deployment: %#v", got)
	}

	prober.jobs = []service.DeploymentJob{job("job", core.DeploymentDeploying)}
	for range 2 {
		if err := svc.PollDeployments(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if got := deploymentByID(t, store, deployment.RequestID); got.LatestStatus != core.DeploymentDeploying || len(got.Events) != 2 {
		t.Fatalf("deployment = %#v, want a single DEPLOYING event", got)
	}
}

func TestPollDeploymentsReadsTheJobListOncePerTick(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{}
	svc, store, _ := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	createQueuedDeployment(t, store, "9b2f0a51-8a9c-4da3-9fed-d20a41cc7317", 3, "A01", now.Add(-time.Minute))
	createQueuedDeployment(t, store, "c1d3e5f7-8a9c-4da3-9fed-d20a41cc7317", 2, "B01", now.Add(-time.Minute))
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if prober.calls != 2 {
		t.Fatalf("ListJobs called %d times over two ticks, want 2", prober.calls)
	}
}

func TestPollDeploymentsSkipsTheGatewayWhenNothingIsActive(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	message := "rstate job ended with state エラー"
	failed := job("job", core.DeploymentFailed)
	failed.Message = &message
	prober := &fakeProber{jobs: []service.DeploymentJob{job("job", core.DeploymentDeploying)}}
	svc, store, deployment := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	prober.jobs = []service.DeploymentJob{failed}
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	got := deploymentByID(t, store, deployment.RequestID)
	last := got.Events[len(got.Events)-1]
	if got.LatestStatus != core.DeploymentFailed || last.Message == nil || *last.Message != message {
		t.Fatalf("deployment = %#v, want FAILED with message", got)
	}

	calls := prober.calls
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if prober.calls != calls {
		t.Fatal("ListJobs was called with no active deployment")
	}
}

func TestPollDeploymentsIsNoOpForCallbackGateways(t *testing.T) {
	store := memory.NewCompetitionStore()
	svc := service.New(store, memory.NewSessionStore(), service.Config{})
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestSyncPolledDeploymentFailsAQueuedRequestWithoutAJob(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{jobs: []service.DeploymentJob{job("old", core.DeploymentCompleted)}}
	svc, store, deployment := newPollingService(t, now, prober)

	err := svc.SyncPolledDeployment(t.Context(), prober, deployment)
	var domainErr *core.Error
	if !errors.As(err, &domainErr) || domainErr.Status != http.StatusNotFound {
		t.Fatalf("error = %v, want 404", err)
	}
	if got := deploymentByID(t, store, deployment.RequestID); got.LatestStatus != core.DeploymentFailed {
		t.Fatalf("deployment = %#v, want FAILED", got)
	}

	calls := prober.calls
	if err := svc.SyncPolledDeployment(t.Context(), prober, deploymentByID(t, store, deployment.RequestID)); err != nil {
		t.Fatalf("sync of a finished deployment: %v", err)
	}
	if prober.calls != calls {
		t.Fatal("a finished deployment was looked up again")
	}
}
