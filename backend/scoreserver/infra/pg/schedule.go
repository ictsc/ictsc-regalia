package pg

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type scheduleRow struct {
	ID      uuid.UUID          `db:"id"`
	Phase   domain.StringPhase `db:"phase"`
	StartAt time.Time          `db:"start_at"`
	EndAt   time.Time          `db:"end_at"`
}

var _ domain.ScheduleReader = (*repo)(nil)

func (r *repo) GetSchedule(ctx context.Context) ([]*domain.ScheduleData, error) {
	var rows []scheduleRow
	if err := sqlx.SelectContext(ctx, r.ext, &rows, `
        SELECT id, phase, start_at, end_at FROM schedules ORDER BY start_at ASC
    `); err != nil {
		return nil, errors.Wrap(err, "failed to get schedules")
	}

	schedules := make([]*domain.ScheduleData, 0, len(rows))
	for _, row := range rows {
		schedules = append(schedules, &domain.ScheduleData{
			ID:          row.ID,
			StringPhase: row.Phase,
			StartAt:     row.StartAt,
			EndAt:       row.EndAt,
		})
	}
	return schedules, nil
}

func (r *RepositoryTx) SaveSchedule(ctx context.Context, data []*domain.ScheduleData) error {
	if _, err := r.ext.ExecContext(ctx, `DELETE FROM schedules`); err != nil {
		return errors.Wrap(err, "failed to delete schedules")
	}

	if len(data) == 0 {
		return nil
	}

	rows := make([]scheduleRow, 0, len(data))
	for _, d := range data {
		rows = append(rows, scheduleRow{
			ID:      d.ID,
			Phase:   d.StringPhase,
			StartAt: d.StartAt,
			EndAt:   d.EndAt,
		})
	}

	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
        INSERT INTO schedules (id, phase, start_at, end_at)
        VALUES (:id, :phase, :start_at, :end_at)
    `, rows); err != nil {
		return errors.Wrap(err, "failed to insert schedules")
	}
	return nil
}
