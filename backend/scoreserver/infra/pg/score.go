package pg

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

var _ domain.UpdateAnswerScoreEffect = (*repo)(nil)

func (r *repo) UpdatePublicAnswerScore(ctx context.Context, input *domain.UpdateAnswerScoreInput) error {
	return r.updateAnswerScore(ctx, input, true)
}

func (r *repo) UpdatePrivateAnswerScore(ctx context.Context, input *domain.UpdateAnswerScoreInput) error {
	return r.updateAnswerScore(ctx, input, false)
}

func (r *repo) updateAnswerScore(ctx context.Context, input *domain.UpdateAnswerScoreInput, isPublic bool) error {
	visibility := "PRIVATE"
	if isPublic {
		visibility = "PUBLIC"
	}
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO answer_scores (answer_id, visibility, marking_result_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (answer_id, visibility) DO UPDATE SET
			marking_result_id = EXCLUDED.marking_result_id`,
		input.AnswerID, visibility, input.MarkingResultID,
	); err != nil {
		return errors.Wrap(err, "failed to insert answer_scores")
	}
	return nil
}

func (r *repo) UpdatePublicProblemScore(ctx context.Context, input *domain.UpdateProblemScoreInput) error {
	return r.updateProblemScore(ctx, input, true)
}

func (r *repo) UpdatePrivateProblemScore(ctx context.Context, input *domain.UpdateProblemScoreInput) error {
	return r.updateProblemScore(ctx, input, false)
}

func (r *repo) updateProblemScore(ctx context.Context, input *domain.UpdateProblemScoreInput, isPublic bool) error {
	visibility := "PRIVATE"
	if isPublic {
		visibility = "PUBLIC"
	}
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO problem_scores
			(team_id, problem_id, visibility, marking_result_id, updated_at)
		VALUES
			($1, $2, $3, $4, $5)
		ON CONFLICT (team_id, problem_id, visibility) DO UPDATE SET
			marking_result_id = EXCLUDED.marking_result_id,
			updated_at = EXCLUDED.updated_at`,
		input.TeamID, input.ProblemID, visibility, input.MarkingResultID, input.UpdateSubmitAt,
	); err != nil {
		return errors.Wrap(err, "failed to insert problem_scores")
	}
	return nil
}
