package memory

import (
	"context"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func TestActivateContentSetsScheduleFreeze(t *testing.T) {
	ctx := context.Background()
	store := NewCompetitionStore()
	freeze := time.Date(2026, 3, 15, 3, 0, 0, 0, time.UTC)
	snapshot := core.ContentSnapshot{CommitSHA: "aaa", Manifest: core.Manifest{RankingFreezeAt: &freeze}}
	if _, err := store.ActivateContent(ctx, snapshot, core.NoActiveContentCommit); err != nil {
		t.Fatal(err)
	}
	state, _ := store.GetCompetitionState(ctx)
	if state.RankingFreezeAt == nil || !state.RankingFreezeAt.Equal(freeze) || state.UpdatedBy != "content:aaa" {
		t.Fatalf("freeze = %v by %q", state.RankingFreezeAt, state.UpdatedBy)
	}

	// Activating the active commit again repairs a freeze that drifted.
	store.state.RankingFreezeAt = nil
	if _, err := store.ActivateContent(ctx, snapshot, "aaa"); err != nil {
		t.Fatal(err)
	}
	if state, _ = store.GetCompetitionState(ctx); state.RankingFreezeAt == nil || !state.RankingFreezeAt.Equal(freeze) {
		t.Fatalf("same-commit activation did not repair the freeze: %v", state.RankingFreezeAt)
	}

	// A schedule without ranking_freeze_at clears it.
	next := core.ContentSnapshot{CommitSHA: "bbb"}
	if _, err := store.ActivateContent(ctx, next, "aaa"); err != nil {
		t.Fatal(err)
	}
	if state, _ = store.GetCompetitionState(ctx); state.RankingFreezeAt != nil {
		t.Fatalf("freeze was not cleared: %v", state.RankingFreezeAt)
	}
}
