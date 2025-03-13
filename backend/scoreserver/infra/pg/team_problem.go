package pg

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

var listTeamProblemsQuery = `
	SELECT
		problem_id, marked_score, penalty, total_score
	FROM team_problem_scores
	WHERE team_id = $1`

func (r *repo) ListProblemsScoreByTeamID(ctx context.Context, teamID uuid.UUID) ([]*domain.TeamProblemScoreData, error) {
	rows, err := r.ext.QueryxContext(ctx, listTeamProblemsQuery, teamID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list team problems")
	}
	defer func() { _ = rows.Close() }()

	var scores []*domain.TeamProblemScoreData
	for rows.Next() {
		var score domain.TeamProblemScoreData
		if err := rows.StructScan(&score); err != nil {
			return nil, errors.Wrap(err, "failed to scan team problem score row")
		}
		scores = append(scores, &score)
	}
	return scores, nil
}
