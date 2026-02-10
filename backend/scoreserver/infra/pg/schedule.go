package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type scheduleRow struct {
	ID      uuid.UUID      `db:"id"`
	Name    sql.NullString `db:"name"`
	StartAt time.Time      `db:"start_at"`
	EndAt   time.Time      `db:"end_at"`
}

func (r *scheduleRow) data() *domain.ScheduleData {
	name := ""
	if r.Name.Valid {
		name = r.Name.String
	}
	return &domain.ScheduleData{
		ID:      r.ID,
		Name:    name,
		StartAt: r.StartAt,
		EndAt:   r.EndAt,
	}
}

func (r *repo) GetSchedule(ctx context.Context) ([]*domain.ScheduleData, error) {
	var rows []scheduleRow
	if err := sqlx.SelectContext(ctx, r.ext, &rows, `
        SELECT id, name, start_at, end_at
        FROM schedules
        WHERE name IS NOT NULL AND name != ''
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
	// 既存のスケジュールを取得して name でマッピング
	existingSchedules, err := r.GetSchedule(ctx)
	if err != nil {
		return err
	}

	existingByName := make(map[string]uuid.UUID) // name -> id
	existingByID := make(map[uuid.UUID]bool)
	for _, existing := range existingSchedules {
		existingByID[existing.ID] = true
		if existing.Name != "" {
			existingByName[existing.Name] = existing.ID
		}
	}

	// リクエストで指定されたスケジュールを処理
	updatedIDs := make(map[uuid.UUID]bool)
	for _, d := range data {
		if d.Name == "" {
			continue // 名前がない場合はスキップ
		}

		var scheduleID uuid.UUID
		if existingID, exists := existingByName[d.Name]; exists {
			scheduleID = existingID // 既存を更新
		} else {
			scheduleID = uuid.Must(uuid.NewV4()) // 新規作成
		}

		// UPSERT (phaseは互換性のためIN_CONTESTを固定で入れる)
		query := `
            INSERT INTO schedules (id, name, phase, start_at, end_at)
            VALUES ($1, $2, 'IN_CONTEST', $3, $4)
            ON CONFLICT (id) DO UPDATE SET
                name = EXCLUDED.name,
                start_at = EXCLUDED.start_at,
                end_at = EXCLUDED.end_at
        `

		if _, err := r.ext.ExecContext(ctx, query, scheduleID, d.Name, d.StartAt, d.EndAt); err != nil {
			return errors.Wrap(err, "failed to upsert schedule")
		}

		updatedIDs[scheduleID] = true
	}

	// リクエストに含まれなかった既存スケジュールを削除
	for id := range existingByID {
		if !updatedIDs[id] {
			if _, err := r.ext.ExecContext(ctx, `DELETE FROM schedules WHERE id = $1`, id); err != nil {
				return errors.Wrap(err, "failed to delete old schedule")
			}
		}
	}

	return nil
}

// GetScheduleIDsByNamesはスケジュール名からIDを解決
func (r *repo) GetScheduleIDsByNames(ctx context.Context, names []string) (map[string]uuid.UUID, error) {
	if len(names) == 0 {
		return make(map[string]uuid.UUID), nil
	}

	query := `SELECT id, name FROM schedules WHERE name = ANY($1)`
	rows, err := r.ext.QueryxContext(ctx, query, names)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schedule IDs by names")
	}
	defer rows.Close()

	result := make(map[string]uuid.UUID)
	for rows.Next() {
		var id uuid.UUID
		var name sql.NullString
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrap(err, "failed to scan schedule")
		}
		if name.Valid {
			result[name.String] = id
		}
	}

	// スケジュールが存在しない場合はエラー
	for _, name := range names {
		if _, ok := result[name]; !ok {
			return nil, errors.Errorf("schedule not found: %s", name)
		}
	}

	return result, nil
}

func (r *repo) GetScheduleNamesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error) {
	if len(ids) == 0 {
		return make(map[uuid.UUID]string), nil
	}

	query := `SELECT id, name FROM schedules WHERE id = ANY($1)`
	rows, err := r.ext.QueryxContext(ctx, query, ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schedule names by IDs")
	}
	defer rows.Close()

	result := make(map[uuid.UUID]string)
	for rows.Next() {
		var id uuid.UUID
		var name sql.NullString
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrap(err, "failed to scan schedule")
		}

		if name.Valid {
			result[id] = name.String
		}
	}

	return result, nil
}
