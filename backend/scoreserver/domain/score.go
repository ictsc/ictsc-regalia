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
	return max(0, s.MarkedScore()-s.Penalty())
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

func (a *Answer) UpdateScore(ctx context.Context, eff UpdateAnswerScoreEffect, now time.Time) error {
	var errs []error
	if err := a.updatePrivateScore(ctx, eff); err != nil {
		errs = append(errs, err)
	}
	if err := a.updatePublicScore(ctx, eff, now); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (a *Answer) updatePublicScore(ctx context.Context, eff UpdateAnswerScoreEffect, now time.Time) error {
	latestPublicMark, err := a.latestPublicMarkingResult(ctx, eff, now)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
	}
	if err := eff.UpdatePublicAnswerScore(ctx, &UpdateAnswerScoreInput{
		AnswerID:        a.id,
		MarkingResultID: latestPublicMark.id,
	}); err != nil {
		return WrapAsInternal(err, "failed to update answer score")
	}
	return nil
}

func (a *Answer) updatePrivateScore(ctx context.Context, eff UpdateAnswerScoreEffect) error {
	latestMark, err := a.latestMarkingResult(ctx, eff)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
	}
	if err := eff.UpdatePrivateAnswerScore(ctx, &UpdateAnswerScoreInput{
		AnswerID:        a.id,
		MarkingResultID: latestMark.id,
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

func (tp *TeamProblem) UpdateScore(ctx context.Context, eff UpdateProblemScoreEffect) error {
	var errs []error
	if err := tp.updatePrivateScore(ctx, eff); err != nil {
		errs = append(errs, err)
	}
	if err := tp.updatePublicScore(ctx, eff); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (tp *TeamProblem) updatePrivateScore(ctx context.Context, eff UpdateProblemScoreEffect) error {
	answers, err := tp.answersForAdmin(ctx, eff)
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

	if err := eff.UpdatePrivateProblemScore(ctx, updates); err != nil {
		return WrapAsInternal(err, "failed to update problem score")
	}

	return nil
}

func (tp *TeamProblem) updatePublicScore(ctx context.Context, eff UpdateProblemScoreEffect) error {
	answers, err := tp.answersForPublic(ctx, eff)
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

	if err := eff.UpdatePublicProblemScore(ctx, updates); err != nil {
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
	}
	UpdateProblemScoreInput struct {
		TeamID          uuid.UUID
		ProblemID       uuid.UUID
		MarkingResultID uuid.UUID
		UpdateSubmitAt  time.Time
	}
	ScoreWriter interface {
		UpdatePrivateAnswerScore(ctx context.Context, input *UpdateAnswerScoreInput) error
		UpdatePublicAnswerScore(ctx context.Context, input *UpdateAnswerScoreInput) error
		UpdatePrivateProblemScore(ctx context.Context, input *UpdateProblemScoreInput) error
		UpdatePublicProblemScore(ctx context.Context, input *UpdateProblemScoreInput) error
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
