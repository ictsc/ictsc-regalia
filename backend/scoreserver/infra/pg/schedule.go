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
	rows := make([]scheduleRow, 0, len(data))
	for _, schedule := range data {
		rows = append(rows, scheduleRow{
			Name:    schedule.Name,
			StartAt: schedule.StartAt,
			EndAt:   schedule.EndAt,
		})
	}

	if _, err := r.ext.ExecContext(ctx, `LOCK TABLE schedules IN ACCESS EXCLUSIVE MODE`); err != nil {
		return errors.Wrap(err, "failed to lock schedules")
	}

	if _, err := r.ext.ExecContext(ctx, `DELETE FROM schedules`); err != nil {
		return errors.Wrap(err, "failed to delete schedules")
	}

	if len(rows) == 0 {
		return nil
	}

	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO schedules (name, start_at, end_at)
		VALUES (:name, :start_at, :end_at)
	`, rows); err != nil {
		return errors.Wrap(err, "failed to insert schedules")
	}
	return nil
}
