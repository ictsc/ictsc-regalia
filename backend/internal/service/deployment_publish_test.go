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

type failOnceEventBus struct {
	publishCalls int
}

func (b *failOnceEventBus) Publish(context.Context, []byte) error {
	b.publishCalls++
	if b.publishCalls == 1 {
		return errors.New("redis temporarily unavailable")
	}
	return nil
}

func (*failOnceEventBus) Subscribe(context.Context) (service.Subscription, error) {
	return nil, nil
}

func TestDuplicateDeploymentCallbackRepublishesAfterEventBusFailure(t *testing.T) {
	store := memory.NewCompetitionStore()
	svc := service.New(store, memory.NewSessionStore(), service.Config{})
	bus := &failOnceEventBus{}
	svc.Events = bus

	deployment, err := store.CreateDeployment(t.Context(), core.Deployment{
		RequestID:     "74d8dc35-8a9c-4da3-9fed-d20a41cc7317",
		TeamCode:      2,
		ProblemCode:   "P1",
		ContentCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		RequestedAt:   time.Unix(10, 0).UTC(),
		LatestStatus:  core.DeploymentQueued,
	})
	if err != nil {
		t.Fatal(err)
	}
	event := core.DeploymentEvent{
		EventID:    "702aa55d-c9fc-4c58-b46e-17b176bba691",
		OccurredAt: time.Unix(11, 0).UTC(),
		Status:     core.DeploymentDeploying,
	}

	_, duplicate, err := svc.ApplyDeploymentEvent(t.Context(), deployment.TeamCode, deployment.ProblemCode, deployment.Revision, event, false)
	var domainErr *core.Error
	if !errors.As(err, &domainErr) || domainErr.Status != http.StatusInternalServerError || duplicate {
		t.Fatalf("first callback: duplicate=%v error=%v, want non-duplicate 500", duplicate, err)
	}

	updated, duplicate, err := svc.ApplyDeploymentEvent(t.Context(), deployment.TeamCode, deployment.ProblemCode, deployment.Revision, event, false)
	if err != nil || !duplicate {
		t.Fatalf("retry callback: duplicate=%v error=%v, want duplicate success", duplicate, err)
	}
	if bus.publishCalls != 2 {
		t.Fatalf("publish calls = %d, want 2", bus.publishCalls)
	}
	if updated.LatestStatus != core.DeploymentDeploying || len(updated.Events) != 1 {
		t.Fatalf("updated deployment = %#v, want one persisted DEPLOYING event", updated)
	}
}
