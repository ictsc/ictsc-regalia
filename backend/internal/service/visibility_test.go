package service

import (
	"context"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

type visibilityStore struct {
	core.Store
	answers  []core.Answer
	markings []core.MarkingResult
	state    core.CompetitionState
}

func (s *visibilityStore) ListAnswers(_ context.Context, filter core.AnswerFilter) ([]core.Answer, error) {
	result := make([]core.Answer, 0, len(s.answers))
	for _, answer := range s.answers {
		if filter.TeamCode != nil && answer.TeamCode != *filter.TeamCode || filter.Before != nil && answer.SubmittedAt.After(*filter.Before) {
			continue
		}
		result = append(result, answer)
	}
	return result, nil
}

func (s *visibilityStore) ListMarkingResults(context.Context) ([]core.MarkingResult, error) {
	return append([]core.MarkingResult(nil), s.markings...), nil
}

func (s *visibilityStore) GetCompetitionState(context.Context) (core.CompetitionState, error) {
	return s.state, nil
}

func TestContestantAndPublicScoreVisibility(t *testing.T) {
	submitted := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	store := &visibilityStore{
		answers:  []core.Answer{{TeamCode: 2, ProblemCode: "P1", Number: 1, SubmittedAt: submitted, MaxScore: 100}},
		markings: []core.MarkingResult{{ID: "mark-1", TeamCode: 2, ProblemCode: "P1", AnswerNumber: 1, MarkedScore: 80, CreatedAt: submitted.Add(time.Minute)}},
	}
	service := New(store, nil, Config{})
	now := submitted.Add(core.PublishDelay - time.Nanosecond)
	service.Now = func() time.Time { return now }
	teamScores, err := service.ContestantScores(context.Background(), 2)
	if err != nil || len(teamScores[2]) != 0 {
		t.Fatalf("private team scores = %#v, err=%v", teamScores, err)
	}
	now = submitted.Add(core.PublishDelay)
	teamScores, err = service.ContestantScores(context.Background(), 2)
	if err != nil || teamScores[2]["P1"].EffectiveScore != 80 {
		t.Fatalf("delayed team scores = %#v, err=%v", teamScores, err)
	}
	freeze := submitted.Add(10 * time.Minute)
	public, err := service.publicScores(context.Background(), &freeze, false)
	if err != nil || len(public[2]) != 0 {
		t.Fatalf("freeze snapshot leaked delayed score = %#v, err=%v", public, err)
	}
	now = submitted.Add(time.Minute)
	public, err = service.publicScores(context.Background(), nil, true)
	if err != nil || public[2]["P1"].EffectiveScore != 80 {
		t.Fatalf("final reveal did not bypass delay = %#v, err=%v", public, err)
	}
}
