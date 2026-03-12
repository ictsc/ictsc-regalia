package pg

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

var _ domain.TeamProblemScoreReader = (*repo)(nil)

var (
	teamProblemScoreQueryBase = `
SELECT
	problem_score.team_id, problem_score.problem_id,
	` + scoreColumns.String("score") + `
FROM scores AS score
JOIN problem_scores AS problem_score
	ON score.marking_result_id = problem_score.marking_result_id AND problem_score.visibility = ?`

	teamProblemScoreByTeamIDAndProblemIDQuery = teamProblemScoreQueryBase + `
WHERE problem_score.team_id = ? AND problem_score.problem_id = ?`

	teamProblemScoreByTeamIDQuery = teamProblemScoreQueryBase + `
WHERE problem_score.team_id = ?`

	teamProblemScoresQuery = teamProblemScoreQueryBase
)

type teamProblemScoreRow struct {
	TeamID          uuid.UUID `db:"team_id"`
	ProblemID       uuid.UUID `db:"problem_id"`
	MarkingResultID uuid.UUID `db:"marking_result_id"`
	scoreRow
}

func (r *repo) GetTeamProblemScore(ctx context.Context, visibility domain.ScoreVisibility, teamID, problemID uuid.UUID) (*domain.ScoreData, error) {
	var row teamProblemScoreRow
	if err := sqlx.GetContext(
		ctx, r.ext, &row, r.ext.Rebind(teamProblemScoreByTeamIDAndProblemIDQuery),
		string(visibility), teamID, problemID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("team problem score", nil)
		}
		return nil, errors.Wrap(err, "failed to get team problem score")
	}
	row.scoreRow.MarkingResultID = row.MarkingResultID

	return (*domain.ScoreData)(&row.scoreRow), nil
}

func (r *repo) ListTeamProblemScoresByTeamID(ctx context.Context, visibility domain.ScoreVisibility, teamID uuid.UUID) ([]*domain.TeamProblemScoreData, error) {
	return r.listTeamProblemScores(ctx, teamProblemScoreByTeamIDQuery, visibility, teamID)
}

func (r *repo) ListTeamProblemScores(ctx context.Context, visibility domain.ScoreVisibility) ([]*domain.TeamProblemScoreData, error) {
	return r.listTeamProblemScores(ctx, teamProblemScoresQuery, visibility)
}

func (r *repo) listTeamProblemScores(
	ctx context.Context, query string, visibility domain.ScoreVisibility, args ...any,
) ([]*domain.TeamProblemScoreData, error) {
	rows := make([]teamProblemScoreRow, 0)
	queryArgs := []any{string(visibility)}
	queryArgs = append(queryArgs, args...)
	if err := sqlx.SelectContext(
		ctx, r.ext, &rows, r.ext.Rebind(query), queryArgs...,
	); err != nil {
		return nil, errors.Wrap(err, "failed to list team problem scores")
	}

	teamProblemScores := make([]*domain.TeamProblemScoreData, 0, len(rows))
	for _, row := range rows {
		row.scoreRow.MarkingResultID = row.MarkingResultID
		teamProblemScores = append(teamProblemScores, &domain.TeamProblemScoreData{
			TeamID:    row.TeamID,
			ProblemID: row.ProblemID,
			Score:     domain.ScoreData(row.scoreRow),
		})
	}

	return teamProblemScores, nil
}

var _ domain.UpdateAnswerScoreEffect = (*repo)(nil)

func (r *repo) UpdateAnswerScore(ctx context.Context, input *domain.UpdateAnswerScoreInput) error {
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO answer_scores (answer_id, visibility, marking_result_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (answer_id, visibility) DO UPDATE SET
			marking_result_id = EXCLUDED.marking_result_id`,
		input.AnswerID, string(input.Visibility), input.MarkingResultID,
	); err != nil {
		return errors.Wrap(err, "failed to insert answer_scores")
	}
	return nil
}

func (r *repo) UpdateProblemScore(ctx context.Context, input *domain.UpdateProblemScoreInput) error {
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO problem_scores
			(team_id, problem_id, visibility, marking_result_id, updated_at)
		VALUES
			($1, $2, $3, $4, $5)
		ON CONFLICT (team_id, problem_id, visibility) DO UPDATE SET
			marking_result_id = EXCLUDED.marking_result_id,
			updated_at = EXCLUDED.updated_at`,
		input.TeamID, input.ProblemID, string(input.Visibility), input.MarkingResultID, input.UpdateSubmitAt,
	); err != nil {
		return errors.Wrap(err, "failed to insert problem_scores")
	}
	return nil
}
