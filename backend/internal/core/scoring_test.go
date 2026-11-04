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

func TestSelectLatestScoresIgnoresOlderMarkedAnswers(t *testing.T) {
	now := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	answers := []Answer{
		{TeamCode: 2, ProblemCode: "A", Number: 1, MaxScore: 100, SubmittedAt: now},
		{TeamCode: 2, ProblemCode: "A", Number: 2, MaxScore: 100, SubmittedAt: now.Add(time.Hour)},
		{TeamCode: 3, ProblemCode: "A", Number: 1, MaxScore: 100, SubmittedAt: now},
		{TeamCode: 3, ProblemCode: "A", Number: 2, MaxScore: 100, SubmittedAt: now.Add(time.Hour)},
	}
	markings := []MarkingResult{
		{ID: "1", TeamCode: 2, ProblemCode: "A", AnswerNumber: 1, MarkedScore: 100, CreatedAt: now},
		{ID: "2", TeamCode: 2, ProblemCode: "A", AnswerNumber: 2, MarkedScore: 40, CreatedAt: now},
		{ID: "3", TeamCode: 3, ProblemCode: "A", AnswerNumber: 1, MarkedScore: 100, CreatedAt: now},
	}
	got := SelectLatestScores(answers, markings)
	if score := got[2]["A"]; score.AnswerNumber != 2 || score.EffectiveScore != 40 {
		t.Fatalf("team 2 got answer %d score %d, want answer 2 score 40", score.AnswerNumber, score.EffectiveScore)
	}
	if _, ok := got[3]["A"]; ok {
		t.Fatalf("team 3 has an unmarked latest answer but got %#v", got[3]["A"])
	}
}

func TestShiftSectionsToPhaseKeepsLayout(t *testing.T) {
	base := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	sections := []Section{
		{Slug: "day1-am", Beginning: base, Ending: base.Add(3 * time.Hour)},
		{Slug: "day1-pm", Beginning: base.Add(4 * time.Hour), Ending: base.Add(7 * time.Hour)},
	}
	now := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	for _, phase := range []ContestPhase{ContestPhaseBefore, ContestPhaseDuring, ContestPhaseAfter} {
		shifted := ShiftSectionsToPhase(sections, phase, now)
		if got := ContestPhaseAt(shifted, now); got != phase {
			t.Fatalf("%s: phase after shift = %s", phase, got)
		}
		if shifted[1].Beginning.Sub(shifted[0].Ending) != time.Hour || shifted[1].Ending.Sub(shifted[1].Beginning) != 3*time.Hour {
			t.Fatalf("%s: layout changed: %#v", phase, shifted)
		}
	}
	if got := ContestPhaseAt(sections, base.Add(3*time.Hour+time.Minute)); got != ContestPhaseAfter {
		t.Fatalf("break between sections = %s, want AFTER", got)
	}
}

func TestContestantVisibleImagesFollowContentVisibility(t *testing.T) {
	now := time.Date(2026, 12, 12, 4, 0, 0, 0, time.UTC)
	snapshot := ContentSnapshot{Manifest: Manifest{
		Sections: []Section{
			{Slug: "open", Beginning: now.Add(-time.Hour), Ending: now.Add(time.Hour)},
			{Slug: "later", Beginning: now.Add(2 * time.Hour), Ending: now.Add(3 * time.Hour)},
		},
		Problems: []Problem{
			{Code: "A", SectionSlug: "open", Body: "![a](/images/open.png)", Explanation: "![s](/images/secret.png)"},
			{Code: "B", SectionSlug: "later", Body: "![b](/images/later.png)"},
		},
		Announcements: []Announcement{
			{Slug: "now", EffectiveFrom: now.Add(-time.Minute), Markdown: "![n](/images/notice.png)"},
			{Slug: "soon", EffectiveFrom: now.Add(time.Minute), Markdown: "![s](/images/soon.png)"},
		},
		RuleMarkdown: "![r](/images/rule.png)",
	}}
	visible := ContestantVisibleImages(snapshot, now)
	for _, name := range []string{"open.png", "notice.png", "rule.png"} {
		if _, ok := visible[name]; !ok {
			t.Errorf("%s should be visible", name)
		}
	}
	for _, name := range []string{"secret.png", "later.png", "soon.png"} {
		if _, ok := visible[name]; ok {
			t.Errorf("%s should be hidden", name)
		}
	}
}
