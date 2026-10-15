package service

import (
	"context"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// ContestantScores returns only scores visible to members of one team. A
// marking remains private until the answer publication delay has elapsed;
// final reveal intentionally bypasses that delay.
func (s *Service) ContestantScores(ctx context.Context, teamCode int64) (map[int64]map[string]core.Score, error) {
	answers, err := s.Store.ListAnswers(ctx, core.AnswerFilter{TeamCode: &teamCode})
	if err != nil {
		return nil, err
	}
	markings, err := s.Store.ListMarkingResults(ctx)
	if err != nil {
		return nil, err
	}
	state, err := s.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	latest := core.LatestMarkings(markings)
	visible := make([]core.MarkingResult, 0, len(answers))
	for _, answer := range answers {
		marking, ok := core.LatestMarkingForAnswer(latest, answer.TeamCode, answer.ProblemCode, answer.Number)
		if !ok {
			continue
		}
		if state.FinalRevealedAt == nil && s.Now().Before(answer.SubmittedAt.Add(core.PublishDelay)) {
			continue
		}
		visible = append(visible, marking)
	}
	return core.SelectBestScores(answers, visible), nil
}

// publicScores evaluates the public view at visibleAt. When visibleAt is the
// freeze instant, both answers and markings are bounded to that instant.
func (s *Service) publicScores(ctx context.Context, visibleAt *time.Time, final bool) (map[int64]map[string]core.Score, error) {
	answers, err := s.Store.ListAnswers(ctx, core.AnswerFilter{Before: visibleAt})
	if err != nil {
		return nil, err
	}
	markings, err := s.Store.ListMarkingResults(ctx)
	if err != nil {
		return nil, err
	}
	instant := s.Now()
	if visibleAt != nil {
		instant = *visibleAt
		filtered := markings[:0]
		for _, marking := range markings {
			if !marking.CreatedAt.After(instant) {
				filtered = append(filtered, marking)
			}
		}
		markings = filtered
	}
	latest := core.LatestMarkings(markings)
	visible := make([]core.MarkingResult, 0, len(answers))
	for _, answer := range answers {
		marking, ok := core.LatestMarkingForAnswer(latest, answer.TeamCode, answer.ProblemCode, answer.Number)
		if !ok || (!final && instant.Before(answer.SubmittedAt.Add(core.PublishDelay))) {
			continue
		}
		visible = append(visible, marking)
	}
	return core.SelectBestScores(answers, visible), nil
}
