package pg

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

var _ domain.RankingReader = (*repo)(nil)

const (
	rankingQuery = `
SELECT
	problem_score.team_id,
	SUM(score.total_score) AS total_score,
	MAX(problem_score.updated_at) AS update_submit_at
FROM problem_scores AS problem_score
INNER JOIN scores AS score ON score.marking_result_id = problem_score.marking_result_id
WHERE problem_score.visibility = $1
GROUP BY problem_score.team_id
ORDER BY total_score DESC, update_submit_at ASC`
)

type rankRow struct {
	TeamID         uuid.UUID `db:"team_id"`
	TotalScore     uint32    `db:"total_score"`
	UpdateSubmitAt time.Time `db:"update_submit_at"`
}

func (r *repo) GetRanking(ctx context.Context, isPublic bool) ([]*domain.RankData, error) {
	visibility := "PRIVATE"
	if isPublic {
		visibility = "PUBLIC"
	}
	rows, err := r.ext.QueryxContext(ctx, rankingQuery, visibility)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select problem_scores")
	}
	defer rows.Close() //nolint:errcheck

	var ranks []*domain.RankData
	for rows.Next() {
		var row rankRow
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		ranks = append(ranks, (*domain.RankData)(&row))
	}
	return ranks, nil
}
