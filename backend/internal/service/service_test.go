package service

import (
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func TestValidateSnapshotRejectsOverlapAndTraversal(t *testing.T) {
	begin := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	threshold, percentage := int32(1), int32(10)
	snapshot := core.ContentSnapshot{
		CommitSHA: "0123456789abcdef0123456789abcdef01234567",
		Manifest: core.Manifest{
			Sections: []core.Section{
				{Slug: "a", Beginning: begin, Ending: begin.Add(2 * time.Hour), ProblemIDs: []string{"A"}},
				{Slug: "b", Beginning: begin.Add(time.Hour), Ending: begin.Add(3 * time.Hour)},
			},
			Problems: []core.Problem{{
				Code: "A", Title: "A", Body: "body", BodyPath: "../secret", MaxScore: 100,
				Redeploy: core.RedeployRule{Type: core.RedeployPercentage, Threshold: &threshold, Percentage: &percentage},
			}},
		},
	}
	if err := ValidateSnapshot(snapshot); err == nil {
		t.Fatal("expected validation error")
	}
}
