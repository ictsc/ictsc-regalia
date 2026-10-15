package memory

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func TestDeploymentEventIdempotencyAndMismatch(t *testing.T) {
	t.Parallel()
	store := NewCompetitionStore()
	ctx := context.Background()
	created, err := store.CreateDeployment(ctx, core.Deployment{
		RequestID: "74d8dc35-8a9c-4da3-9fed-d20a41cc7317", TeamCode: 1, ProblemCode: "A",
		ContentCommit: "0123456789012345678901234567890123456789", RequestedAt: time.Unix(10, 0).UTC(), LatestStatus: core.DeploymentQueued,
	})
	if err != nil {
		t.Fatal(err)
	}
	event := core.DeploymentEvent{EventID: "702aa55d-c9fc-4c58-b46e-17b176bba691", OccurredAt: time.Unix(11, 0).UTC(), Status: core.DeploymentDeploying}
	if _, duplicate, err := store.AppendDeploymentEvent(ctx, created.RequestID, event, false); err != nil || duplicate {
		t.Fatalf("first callback: duplicate=%v err=%v", duplicate, err)
	}
	if _, duplicate, err := store.AppendDeploymentEvent(ctx, created.RequestID, event, false); err != nil || !duplicate {
		t.Fatalf("exact retry: duplicate=%v err=%v", duplicate, err)
	}
	mismatch := event
	mismatch.Status = core.DeploymentCompleted
	_, _, err = store.AppendDeploymentEvent(ctx, created.RequestID, mismatch, false)
	var domainErr *core.Error
	if !errors.As(err, &domainErr) || domainErr.Status != http.StatusConflict || domainErr.Code != "duplicate_event_mismatch" {
		t.Fatalf("mismatch error = %#v", err)
	}
}

func TestOnlyOneDeploymentCanBeActive(t *testing.T) {
	t.Parallel()
	store := NewCompetitionStore()
	ctx := context.Background()
	base := core.Deployment{TeamCode: 1, ProblemCode: "A", ContentCommit: "0123456789012345678901234567890123456789", RequestedAt: time.Now(), LatestStatus: core.DeploymentQueued}
	base.RequestID = "74d8dc35-8a9c-4da3-9fed-d20a41cc7317"
	if _, err := store.CreateDeployment(ctx, base); err != nil {
		t.Fatal(err)
	}
	base.RequestID = "8e9415b5-21d3-47bd-8b5c-b3bfb0a783eb"
	_, err := store.CreateDeployment(ctx, base)
	var domainErr *core.Error
	if !errors.As(err, &domainErr) || domainErr.Code != "deployment_in_progress" {
		t.Fatalf("second create error = %#v", err)
	}
}
