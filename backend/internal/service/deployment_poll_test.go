package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

type fakeProber struct {
	observation service.DeploymentObservation
	err         error
	calls       int
}

func (*fakeProber) Queue(context.Context, service.DeploymentRequest) error { return nil }

func (*fakeProber) Status(context.Context, int64, string) (core.DeploymentEvent, error) {
	return core.DeploymentEvent{}, nil
}

func (p *fakeProber) Observe(context.Context, int64, string) (service.DeploymentObservation, error) {
	p.calls++
	return p.observation, p.err
}

type notFoundError struct{}

func (notFoundError) Error() string                  { return "no job" }
func (notFoundError) DeploymentStatusNotFound() bool { return true }

func newPollingService(t *testing.T, now time.Time, prober service.DeploymentGateway) (*service.Service, *memory.CompetitionStore, core.Deployment) {
	t.Helper()
	store := memory.NewCompetitionStore()
	svc := service.New(store, memory.NewSessionStore(), service.Config{})
	svc.Now = func() time.Time { return now }
	svc.Deployments = prober
	requestedAt := now.Add(-time.Minute)
	deployment, err := store.CreateDeployment(t.Context(), core.Deployment{
		RequestID: "74d8dc35-8a9c-4da3-9fed-d20a41cc7317", TeamCode: 2, ProblemCode: "A01",
		ContentCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", RequestedAt: requestedAt, LatestStatus: core.DeploymentQueued,
		Events: []core.DeploymentEvent{{EventID: "702aa55d-c9fc-4c58-b46e-17b176bba691", OccurredAt: requestedAt, Status: core.DeploymentQueued}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return svc, store, deployment
}

func latestDeployment(t *testing.T, store *memory.CompetitionStore) core.Deployment {
	t.Helper()
	deployments, err := store.ListDeployments(t.Context(), core.DeploymentFilter{})
	if err != nil || len(deployments) != 1 {
		t.Fatalf("deployments = %#v, err = %v", deployments, err)
	}
	return deployments[0]
}

func TestPollDeploymentsCompletesQueuedDeploymentThroughDeploying(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{observation: service.DeploymentObservation{JobID: "job", CreatedAt: now.Add(-30 * time.Second), Status: core.DeploymentCompleted}}
	svc, store, _ := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	deployment := latestDeployment(t, store)
	if deployment.LatestStatus != core.DeploymentCompleted || len(deployment.Events) != 3 {
		t.Fatalf("deployment = %#v, want QUEUED -> DEPLOYING -> COMPLETED", deployment)
	}
	if !deployment.Events[2].OccurredAt.After(deployment.Events[1].OccurredAt) {
		t.Fatalf("events are not strictly ordered: %#v", deployment.Events)
	}

	// A finished deployment is no longer polled.
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if prober.calls != 1 {
		t.Fatalf("Observe called %d times, want 1", prober.calls)
	}
}

func TestPollDeploymentsKeepsWaitingAndRepeatedStatesQuiet(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{observation: service.DeploymentObservation{CreatedAt: now, Status: core.DeploymentQueued}}
	svc, store, _ := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := latestDeployment(t, store); got.LatestStatus != core.DeploymentQueued || len(got.Events) != 1 {
		t.Fatalf("waiting job changed the deployment: %#v", got)
	}

	prober.observation.Status = core.DeploymentDeploying
	for range 2 {
		if err := svc.PollDeployments(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if got := latestDeployment(t, store); got.LatestStatus != core.DeploymentDeploying || len(got.Events) != 2 {
		t.Fatalf("deployment = %#v, want a single DEPLOYING event", got)
	}
}

func TestPollDeploymentsIgnoresJobsOlderThanTheRequest(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{observation: service.DeploymentObservation{CreatedAt: now.Add(-time.Hour), Status: core.DeploymentCompleted}}
	svc, store, _ := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	if got := latestDeployment(t, store); got.LatestStatus != core.DeploymentQueued {
		t.Fatalf("stale job moved the deployment to %s", got.LatestStatus)
	}
}

func TestPollDeploymentsRecordsFailureMessageAndSkipsMissingJobs(t *testing.T) {
	now := time.Date(2026, 11, 4, 3, 0, 0, 0, time.UTC)
	prober := &fakeProber{err: notFoundError{}}
	svc, store, _ := newPollingService(t, now, prober)

	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatalf("missing job should be skipped: %v", err)
	}

	message := "rstate job ended with state エラー"
	prober.err = nil
	prober.observation = service.DeploymentObservation{CreatedAt: now, Status: core.DeploymentFailed, Message: &message}
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
	got := latestDeployment(t, store)
	last := got.Events[len(got.Events)-1]
	if got.LatestStatus != core.DeploymentFailed || last.Message == nil || *last.Message != message {
		t.Fatalf("deployment = %#v, want FAILED with message", got)
	}
}

func TestPollDeploymentsIsNoOpForCallbackGateways(t *testing.T) {
	store := memory.NewCompetitionStore()
	svc := service.New(store, memory.NewSessionStore(), service.Config{})
	if err := svc.PollDeployments(t.Context()); err != nil {
		t.Fatal(err)
	}
}
