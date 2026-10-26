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
	state, err := s.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	if s.ScoresHidden(state) {
		return map[int64]map[string]core.Score{}, nil
	}
	answers, err := s.Store.ListAnswers(ctx, core.AnswerFilter{TeamCode: &teamCode})
	if err != nil {
		return nil, err
	}
	markings, err := s.Store.ListMarkingResults(ctx)
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
	return s.selectScores(answers, visible), nil
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
	return s.selectScores(answers, visible), nil
}

// ScoresHidden reports whether contestants must not see any score yet. In
// secondary qualifier mode every result stays private until the final reveal.
func (s *Service) ScoresHidden(state core.CompetitionState) bool {
	return s.Config.SecondaryQualifierMode && state.FinalRevealedAt == nil
}

// MarkingVisibility extends core.MarkingVisibilityAt with the secondary
// qualifier rule that nothing is released before the final reveal.
func (s *Service) MarkingVisibility(answerSubmittedAt time.Time, state core.CompetitionState) core.Visibility {
	if s.ScoresHidden(state) {
		return core.VisibilityPrivate
	}
	return core.MarkingVisibilityAt(answerSubmittedAt, s.Now(), state.RankingFreezeAt, state.FinalRevealedAt)
}

// selectScores picks the scored answer per team and problem: the latest
// answer in secondary qualifier mode, otherwise the best one.
func (s *Service) selectScores(answers []core.Answer, markings []core.MarkingResult) map[int64]map[string]core.Score {
	if s.Config.SecondaryQualifierMode {
		return core.SelectLatestScores(answers, markings)
	}
	return core.SelectBestScores(answers, markings)
}

// ContestEndAt returns the end of the last section in the active content, or
// nil when no section is scheduled.
func (s *Service) ContestEndAt(ctx context.Context) (*time.Time, error) {
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	var end *time.Time
	for _, section := range snapshot.Manifest.Sections {
		if end == nil || section.Ending.After(*end) {
			ending := section.Ending
			end = &ending
		}
	}
	return end, nil
}

// AnswerGrading reports whether answer is its team's latest for the problem
// and, when it cannot be graded yet, why. The reason is always empty outside
// secondary qualifier mode, where any answer may be graded at any time.
func (s *Service) AnswerGrading(ctx context.Context, answer core.Answer) (bool, string, error) {
	team, problem := answer.TeamCode, answer.ProblemCode
	answers, err := s.Store.ListAnswers(ctx, core.AnswerFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		return false, "", err
	}
	latest := core.IsLatestAnswer(answer, answers)
	if !s.Config.SecondaryQualifierMode {
		return latest, "", nil
	}
	end, err := s.ContestEndAt(ctx)
	if err != nil {
		return false, "", err
	}
	if end != nil && s.Now().Before(*end) {
		return latest, MarkingBlockedContestRunning, nil
	}
	if !latest {
		return latest, MarkingBlockedNotLatest, nil
	}
	return latest, "", nil
}

// CompetitionMode names the scoring mode exposed through the API.
func (s *Service) CompetitionMode() string {
	if s.Config.SecondaryQualifierMode {
		return "SECONDARY_QUALIFIER"
	}
	return "STANDARD"
}

const (
	MarkingBlockedContestRunning = "grading_not_started"
	MarkingBlockedNotLatest      = "answer_not_latest"
)
