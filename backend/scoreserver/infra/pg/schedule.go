package pg

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type schedulePhase domain.Phase

var _ sql.Scanner = (*schedulePhase)(nil)
var _ driver.Valuer = schedulePhase(0)

type scheduleRow struct {
	ID      uuid.UUID     `db:"id"`
	Phase   schedulePhase `db:"phase"`
	StartAt time.Time     `db:"start_at"`
	EndAt   time.Time     `db:"end_at"`
}

func (r *scheduleRow) data() *domain.ScheduleData {
	return &domain.ScheduleData{
		ID:      r.ID,
		Phase:   domain.Phase(r.Phase),
		StartAt: r.StartAt,
		EndAt:   r.EndAt,
	}
}

func (r *repo) GetSchedule(ctx context.Context) ([]*domain.ScheduleData, error) {
	var rows []scheduleRow
	if err := sqlx.SelectContext(ctx, r.ext, &rows, `
        SELECT id, phase, start_at, end_at
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
			Phase:   schedulePhase(d.Phase),
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

func (p *schedulePhase) Scan(src any) error {
	*p = schedulePhase(domain.PhaseUnspecified)

	if src == nil {
		return nil
	}
	v, ok := src.(string)
	if !ok {
		return nil
	}

	switch v {
	case "OUT_OF_CONTEST":
		*p = schedulePhase(domain.PhaseOutOfContest)
	case "IN_CONTEST":
		*p = schedulePhase(domain.PhaseInContest)
	case "BREAK":
		*p = schedulePhase(domain.PhaseBreak)
	case "AFTER_CONTEST":
		*p = schedulePhase(domain.PhaseAfterContest)
	case "UNSPECIFIED":
		// 何もしなくても既に Unspecified
	}
	return nil
}

func (p schedulePhase) Value() (driver.Value, error) {
	phase := domain.Phase(p)
	switch phase {
	case domain.PhaseOutOfContest:
		return "OUT_OF_CONTEST", nil
	case domain.PhaseInContest:
		return "IN_CONTEST", nil
	case domain.PhaseBreak:
		return "BREAK", nil
	case domain.PhaseAfterContest:
		return "AFTER_CONTEST", nil
	case domain.PhaseUnspecified:
		fallthrough
	default:
		return nil, errors.New("unknown phase")
	}
}
