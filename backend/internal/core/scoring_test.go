package core

import (
	"testing"
	"time"
)

func TestDeploymentPenaltyUsesFloorBeforeMultiplication(t *testing.T) {
	threshold, percentage := int32(1), int32(17)
	got := DeploymentPenalty(99, RedeployRule{Type: RedeployPercentage, Threshold: &threshold, Percentage: &percentage}, 4)
	if got != 48 { // floor(99 * 17 / 100) * (4 - 1)
		t.Fatalf("got %d, want 48", got)
	}
}

func TestSelectBestScoreUsesLatestMarkThenBestAnswer(t *testing.T) {
	now := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	answers := []Answer{
		{TeamCode: 2, ProblemCode: "A", Number: 1, MaxScore: 100, SubmittedAt: now},
		{TeamCode: 2, ProblemCode: "A", Number: 2, MaxScore: 100, SubmittedAt: now.Add(time.Hour)},
	}
	markings := []MarkingResult{
		{ID: "1", TeamCode: 2, ProblemCode: "A", AnswerNumber: 1, MarkedScore: 90, CreatedAt: now},
		{ID: "2", TeamCode: 2, ProblemCode: "A", AnswerNumber: 1, MarkedScore: 70, CreatedAt: now.Add(time.Minute)},
		{ID: "3", TeamCode: 2, ProblemCode: "A", AnswerNumber: 2, MarkedScore: 80, CreatedAt: now},
	}
	got := SelectBestScores(answers, markings)[2]["A"]
	if got.AnswerNumber != 2 || got.EffectiveScore != 80 {
		t.Fatalf("got answer %d score %d, want answer 2 score 80", got.AnswerNumber, got.EffectiveScore)
	}
}

func TestBuildRankingDenseTie(t *testing.T) {
	now := time.Now().UTC()
	teams := []Team{{Code: 2}, {Code: 3}, {Code: 4}}
	selected := map[int64]map[string]Score{
		2: {"A": {EffectiveScore: 10, SubmittedAt: now}},
		3: {"A": {EffectiveScore: 10, SubmittedAt: now}},
		4: {"A": {EffectiveScore: 5, SubmittedAt: now}},
	}
	ranking := BuildRanking(teams, map[string]struct{}{"A": {}}, selected)
	if ranking[0].Rank != 1 || ranking[1].Rank != 1 || ranking[2].Rank != 2 {
		t.Fatalf("unexpected ranks: %d %d %d", ranking[0].Rank, ranking[1].Rank, ranking[2].Rank)
	}
}

func TestDeploymentTransitions(t *testing.T) {
	if !CanTransitionDeployment(DeploymentQueued, DeploymentDeploying, false) {
		t.Fatal("queued should transition to deploying")
	}
	if CanTransitionDeployment(DeploymentCompleted, DeploymentFailed, false) {
		t.Fatal("completed must not transition to failed")
	}
}
