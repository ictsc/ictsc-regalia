package memory

import (
	"context"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func TestAnswerWorkflowTransitions(t *testing.T) {
	ctx := context.Background()
	store := NewCompetitionStore()
	answer, _, err := store.SubmitAnswer(ctx, core.Answer{TeamCode: 12, ProblemCode: "A01", SubmittedAt: time.Now()}, 0)
	if err != nil {
		t.Fatal(err)
	}
	workflow, err := store.GetAnswerWorkflow(ctx, 12, "A01", answer.Number)
	if err != nil || workflow.Status != core.AnswerWaiting || workflow.Revision != 0 {
		t.Fatalf("initial workflow: %#v, %v", workflow, err)
	}
	completed := core.AnswerCompleted
	if _, err := store.UpdateAnswerWorkflow(ctx, 12, "A01", answer.Number, core.WorkflowUpdate{Status: &completed}); err == nil {
		t.Fatal("unmarked answer completed")
	}
	inProgress := core.AnswerInProgress
	workflow, err = store.UpdateAnswerWorkflow(ctx, 12, "A01", answer.Number, core.WorkflowUpdate{Status: &inProgress, Assignment: "CLAIM_SELF", ActorDiscordID: "123"})
	if err != nil || workflow.Status != inProgress || workflow.ClaimedBy != "123" || workflow.Revision != 1 {
		t.Fatalf("claimed workflow: %#v, %v", workflow, err)
	}
	if _, err := store.UpdateAnswerWorkflow(ctx, 12, "A01", answer.Number, core.WorkflowUpdate{ExpectedRevision: 0, Assignment: "RESET_TO_DEFAULT"}); err == nil {
		t.Fatal("stale update accepted")
	}
	_, err = store.CreateMarkingResult(ctx, core.MarkingResult{TeamCode: 12, ProblemCode: "A01", AnswerNumber: answer.Number})
	if err != nil {
		t.Fatal(err)
	}
	workflow, _ = store.GetAnswerWorkflow(ctx, 12, "A01", answer.Number)
	if workflow.Status != core.AnswerCompleted || workflow.Revision != 2 || workflow.ClaimedBy != "123" {
		t.Fatalf("marked workflow: %#v", workflow)
	}
	workflow, err = store.UpdateAnswerWorkflow(ctx, 12, "A01", answer.Number, core.WorkflowUpdate{ExpectedRevision: 2, Status: &inProgress, Assignment: "RESET_TO_DEFAULT"})
	if err != nil || workflow.Status != inProgress || workflow.ClaimedBy != "" {
		t.Fatalf("reopened workflow: %#v, %v", workflow, err)
	}
}
