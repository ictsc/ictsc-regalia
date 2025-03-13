package pg

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

var getRankingQuery = `
	WITH best_scores AS (
		SELECT
			team_id,
			total_score,
			created_at
		FROM team_problem_scores
	)
	SELECT
		t.name AS team_name,
		SUM(bs.total_score) AS score,
		MAX(bs.created_at) AS submitted_at
	FROM best_scores bs
	INNER JOIN teams t ON bs.team_id = t.id
	GROUP BY t.id, t.name
	ORDER BY score DESC, submitted_at ASC;
`

func (r *repo) GetRanking(ctx context.Context) ([]*domain.RankingData, error) {
	rows, err := r.ext.QueryxContext(ctx, getRankingQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ranking")
	}
	defer func() { _ = rows.Close() }()

	var ranking []*domain.RankingData
	for rows.Next() {
		var rank domain.RankingData
		if err := rows.StructScan(&rank); err != nil {
			return nil, errors.Wrap(err, "failed to scan ranking row")
		}
		ranking = append(ranking, &rank)
	}

	for i, rd := range ranking {
		rd.Rank = int64(i + 1)
	}
	return ranking, nil
}
