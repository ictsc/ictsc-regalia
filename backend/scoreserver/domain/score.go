package domain

import (
	"context"
	"slices"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type Score struct {
	markingResultID uuid.UUID
	max             uint32
	marked          uint32
	penalty         uint32
}

func (s *Score) MarkedScore() uint32 {
	return s.marked
}

func (s *Score) MaxScore() uint32 {
	return s.max
}

func (s *Score) Penalty() uint32 {
	return s.penalty
}

func (s *Score) TotalScore() uint32 {
	return uint32(max(0, int64(s.MarkedScore())-int64(s.Penalty()))) //nolint:gosec
}

func (s *Score) MarkingResult(ctx context.Context, eff MarkingResultReader) (*MarkingResult, error) {
	marks, err := ListAllMarkingResults(ctx, eff)
	if err != nil {
		return nil, err
	}
	idx := slices.IndexFunc(marks, func(m *MarkingResult) bool {
		return m.id == s.markingResultID
	})
	if idx < 0 {
		return nil, NewNotFoundError("marking result", nil)
	}
	return marks[idx], nil
}

type UpdateAnswerScoreEffect interface {
	MarkingResultReader
	ScoreWriter
}

type ScoreUpdatePolicy struct {
	UpdatePrivate         bool
	UpdateTeam            bool
	UpdatePublic          bool
	BypassVisibilityDelay bool
}

func NewScoreUpdatePolicy(
	mode ScoreUpdateMode,
	inContest bool,
	rankingFrozen bool,
) ScoreUpdatePolicy {
	switch mode {
	case ScoreUpdateModeRevealFinal:
		return ScoreUpdatePolicy{
			UpdatePrivate:         true,
			UpdateTeam:            true,
			UpdatePublic:          true,
			BypassVisibilityDelay: true,
		}
	case ScoreUpdateModeNormal:
		fallthrough
	default:
		return ScoreUpdatePolicy{
			UpdatePrivate: true,
			UpdateTeam:    inContest,
			UpdatePublic:  inContest && !rankingFrozen,
		}
	}
}

func (a *Answer) UpdateScore(
	ctx context.Context,
	eff UpdateAnswerScoreEffect,
	now time.Time,
	policy ScoreUpdatePolicy,
) error {
	var errs []error
	if policy.UpdatePrivate {
		if err := a.updateScore(ctx, eff, ScoreVisibilityPrivate, now, true); err != nil {
			errs = append(errs, err)
		}
	}
	if policy.UpdateTeam {
		if err := a.updateScore(ctx, eff, ScoreVisibilityTeam, now, policy.BypassVisibilityDelay); err != nil {
			errs = append(errs, err)
		}
	}
	if policy.UpdatePublic {
		if err := a.updateScore(ctx, eff, ScoreVisibilityPublic, now, policy.BypassVisibilityDelay); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (a *Answer) updateScore(
	ctx context.Context,
	eff UpdateAnswerScoreEffect,
	visibility ScoreVisibility,
	now time.Time,
	bypassDelay bool,
) error {
	latestMark, err := a.latestMarkingResultForVisibility(ctx, eff, now, visibility, bypassDelay)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return err
	}
	if err := eff.UpdateAnswerScore(ctx, &UpdateAnswerScoreInput{
		AnswerID:        a.id,
		MarkingResultID: latestMark.id,
		Visibility:      visibility,
	}); err != nil {
		return WrapAsInternal(err, "failed to update answer score")
	}
	return nil
}

type UpdateProblemScoreEffect interface {
	AnswerReader
	MarkingResultReader
	ScoreWriter
}

func (tp *TeamProblem) UpdateScore(
	ctx context.Context,
	eff UpdateProblemScoreEffect,
	policy ScoreUpdatePolicy,
) error {
	var errs []error
	if policy.UpdatePrivate {
		if err := tp.updateScore(ctx, eff, ScoreVisibilityPrivate); err != nil {
			errs = append(errs, err)
		}
	}
	if policy.UpdateTeam {
		if err := tp.updateScore(ctx, eff, ScoreVisibilityTeam); err != nil {
			errs = append(errs, err)
		}
	}
	if policy.UpdatePublic {
		if err := tp.updateScore(ctx, eff, ScoreVisibilityPublic); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (tp *TeamProblem) updateScore(
	ctx context.Context,
	eff UpdateProblemScoreEffect,
	visibility ScoreVisibility,
) error {
	var (
		answers []*Answer
		err     error
	)
	switch visibility {
	case ScoreVisibilityPrivate:
		answers, err = tp.answersForAdmin(ctx, eff)
	case ScoreVisibilityTeam:
		answers, err = tp.answersForTeam(ctx, eff)
	case ScoreVisibilityPublic:
		answers, err = tp.answersForPublic(ctx, eff)
	default:
		return NewInvalidArgumentError("unknown score visibility", nil)
	}
	if err != nil {
		return err
	}

	updates, err := tp.problemUpdate(answers)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return err
	}

	updates.Visibility = visibility
	if err := eff.UpdateProblemScore(ctx, updates); err != nil {
		return WrapAsInternal(err, "failed to update problem score")
	}

	return nil
}

func (tp *TeamProblem) problemUpdate(answers []*Answer) (*UpdateProblemScoreInput, error) {
	slices.SortFunc(answers, func(a, b *Answer) int {
		return a.CreatedAt().Compare(b.CreatedAt())
	})

	var scoredAnswer *Answer
	for _, answer := range answers {
		score := answer.Score()
		if score == nil {
			continue // non-scored answer
		}
		if scoredAnswer == nil || answer.Score().TotalScore() > scoredAnswer.Score().TotalScore() {
			scoredAnswer = answer
		}
	}
	if scoredAnswer == nil {
		return nil, NewNotFoundError("scored answer", nil)
	}

	return &UpdateProblemScoreInput{
		TeamID:          uuid.UUID(tp.Team().teamID),
		ProblemID:       uuid.UUID(tp.problemID),
		MarkingResultID: scoredAnswer.Score().markingResultID,
		UpdateSubmitAt:  scoredAnswer.CreatedAt(),
	}, nil
}

type (
	ScoreData struct {
		MarkingResultID uuid.UUID `json:"marking_result_id"`
		MarkedScore     uint32    `json:"marked_score"`
		Penalty         uint32    `json:"penalty"`
	}

	UpdateAnswerScoreInput struct {
		AnswerID        uuid.UUID
		MarkingResultID uuid.UUID
		Visibility      ScoreVisibility
	}
	UpdateProblemScoreInput struct {
		TeamID          uuid.UUID
		ProblemID       uuid.UUID
		MarkingResultID uuid.UUID
		UpdateSubmitAt  time.Time
		Visibility      ScoreVisibility
	}
	ScoreWriter interface {
		UpdateAnswerScore(ctx context.Context, input *UpdateAnswerScoreInput) error
		UpdateProblemScore(ctx context.Context, input *UpdateProblemScoreInput) error
	}
)

func (s *Score) Data() *ScoreData {
	return &ScoreData{
		MarkingResultID: s.markingResultID,
		MarkedScore:     s.marked,
		Penalty:         s.penalty,
	}
}

func (s *ScoreData) parse(problem *Problem) (*Score, error) {
	if s.MarkedScore > problem.MaxScore() {
		return nil, NewInvalidArgumentError("marked score is over max score", nil)
	}
	return &Score{
		markingResultID: s.MarkingResultID,
		max:             problem.maxScore,
		marked:          s.MarkedScore,
		penalty:         s.Penalty,
	}, nil
}
