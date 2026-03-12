package pg

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

var _ domain.ScoreVisibilitySettingsReader = (*repo)(nil)
var _ domain.ScoreVisibilitySettingsWriter = (*RepositoryTx)(nil)

func (r *repo) GetScoreVisibilitySettings(ctx context.Context) (*domain.ScoreVisibilitySettingsData, error) {
	var row struct {
		RankingFreezeAt sql.NullTime `db:"ranking_freeze_at"`
	}
	if err := sqlx.GetContext(ctx, r.ext, &row, `
		SELECT ranking_freeze_at
		FROM score_visibility_settings
		WHERE id = TRUE
		LIMIT 1`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.ScoreVisibilitySettingsData{}, nil
		}
		return nil, errors.Wrap(err, "failed to get score_visibility_settings")
	}

	data := &domain.ScoreVisibilitySettingsData{}
	if row.RankingFreezeAt.Valid {
		ts := row.RankingFreezeAt.Time
		data.RankingFreezeAt = &ts
	}
	return data, nil
}

func (r *RepositoryTx) SaveScoreVisibilitySettings(ctx context.Context, data *domain.ScoreVisibilitySettingsData) error {
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO score_visibility_settings (id, ranking_freeze_at)
		VALUES (TRUE, $1)
		ON CONFLICT (id) DO UPDATE SET
			ranking_freeze_at = EXCLUDED.ranking_freeze_at`,
		data.RankingFreezeAt,
	); err != nil {
		return errors.Wrap(err, "failed to save score_visibility_settings")
	}
	return nil
}
