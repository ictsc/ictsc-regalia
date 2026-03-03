package pg

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type scheduleRow struct {
	Name    string    `db:"name"`
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
}

func (r *scheduleRow) data() *domain.ScheduleData {
	return &domain.ScheduleData{
		Name:    r.Name,
		StartAt: r.StartAt,
		EndAt:   r.EndAt,
	}
}

func (r *repo) GetSchedule(ctx context.Context) ([]*domain.ScheduleData, error) {
	var rows []scheduleRow
	if err := sqlx.SelectContext(ctx, r.ext, &rows, `
        SELECT name, start_at, end_at
        FROM schedules
        ORDER BY start_at ASC
    `); err != nil {
		return nil, errors.Wrap(err, "failed to get schedules")
	}

	schedules := make([]*domain.ScheduleData, 0, len(rows))
	for _, row := range rows {
		schedules = append(schedules, row.data())
	}
	return schedules, nil
}

func (r *RepositoryTx) SaveSchedule(ctx context.Context, data []*domain.ScheduleData) error {
	// 既存のスケジュール名を取得
	existingSchedules, err := r.GetSchedule(ctx)
	if err != nil {
		return err
	}

	existingNames := make(map[string]bool)
	for _, existing := range existingSchedules {
		existingNames[existing.Name] = true
	}

	// リクエストで指定されたスケジュールを処理
	updatedNames := make(map[string]bool)
	for _, d := range data {
		if d.Name == "" {
			continue // 名前がない場合はスキップ
		}

		query := `
            INSERT INTO schedules (name, start_at, end_at)
            VALUES ($1, $2, $3)
            ON CONFLICT (name) DO UPDATE SET
                start_at = EXCLUDED.start_at,
                end_at = EXCLUDED.end_at
        `

		if _, err := r.ext.ExecContext(ctx, query, d.Name, d.StartAt, d.EndAt); err != nil {
			return errors.Wrap(err, "failed to upsert schedule")
		}

		updatedNames[d.Name] = true
	}

	// リクエストに含まれなかった既存スケジュールを削除
	for name := range existingNames {
		if !updatedNames[name] {
			if _, err := r.ext.ExecContext(ctx, `DELETE FROM schedules WHERE name = $1`, name); err != nil {
				return errors.Wrap(err, "failed to delete old schedule")
			}
		}
	}

	return nil
}
