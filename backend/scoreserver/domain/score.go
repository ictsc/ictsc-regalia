package domain

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
)

type Score struct {
	max     uint32
	marked  uint32
	penalty uint32
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

type (
	ScoreData struct {
		MarkedScore uint32 `json:"marked_score"`
		Penalty     uint32 `json:"penalty"`
	}

	UpdateAnswerScoreInput struct {
		AnswerID        uuid.UUID
		MarkingResultID uuid.UUID
	}
	ScoreWriter interface {
		UpdatePrivateAnswerScore(ctx context.Context, input *UpdateAnswerScoreInput) error
		UpdatePublicAnswerScore(ctx context.Context, input *UpdateAnswerScoreInput) error
	}
)

func (s *Score) Data() *ScoreData {
	return &ScoreData{
		MarkedScore: s.marked,
		Penalty:     s.penalty,
	}
}

func (s *ScoreData) parse(problem *Problem) (*Score, error) {
	if s.MarkedScore > problem.MaxScore() {
		return nil, NewInvalidArgumentError("marked score is over max score", nil)
	}
	return &Score{
		max:     problem.maxScore,
		marked:  s.MarkedScore,
		penalty: s.Penalty,
	}, nil
}
