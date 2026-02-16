package pg_test

import (
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.GetSchedule(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// seed.sql に 4 件のスケジュールが登録されている
	if len(actual) != 4 {
		t.Fatalf("expected 4 schedules, got %d", len(actual))
	}

	slices.SortFunc(actual, func(a, b *domain.ScheduleData) int {
		return a.StartAt.Compare(b.StartAt)
	})

	expectedNames := []string{"day1-am", "day1-pm", "day2-am", "day2-pm"}
	for i, name := range expectedNames {
		if actual[i].Name != name {
			t.Errorf("schedule[%d].Name = %q, want %q", i, actual[i].Name, name)
		}
	}
}

func TestSaveSchedule(t *testing.T) {
	t.Parallel()

	startAt := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC)

	cases := map[string]struct {
		input   []*domain.ScheduleData
		queries []string
	}{
		"single schedule": {
			input: []*domain.ScheduleData{
				{
					Name:    "new-schedule",
					StartAt: startAt,
					EndAt:   endAt,
				},
			},
			queries: []string{
				`SELECT 1
				 FROM schedules
				 WHERE name = 'new-schedule'
				   AND start_at = '2025-02-03 00:00:00'
				   AND end_at = '2025-02-04 00:00:00'`,
			},
		},
		"multiple schedules": {
			input: []*domain.ScheduleData{
				{
					Name:    "sched-1",
					StartAt: startAt,
					EndAt:   endAt,
				},
				{
					Name:    "sched-2",
					StartAt: endAt,
					EndAt:   endAt.Add(24 * time.Hour),
				},
			},
			queries: []string{
				`SELECT 1
				FROM schedules
				WHERE name = 'sched-1'
				AND start_at = '2025-02-03 00:00:00'
				AND end_at = '2025-02-04 00:00:00'`,
				`SELECT 1
				FROM schedules
				WHERE name = 'sched-2'
				AND start_at = '2025-02-04 00:00:00'
				AND end_at = '2025-02-05 00:00:00'`,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)

			err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
				return tx.SaveSchedule(t.Context(), tt.input)
			})
			if err != nil {
				t.Fatalf("failed to save schedule: %+v", err)
			}

			for _, query := range tt.queries {
				var dst any
				if err := db.GetContext(t.Context(), &dst, query); err != nil {
					t.Errorf("query %s: %v", query, err)
				}
			}
		})
	}
}

func TestSaveScheduleUpdate(t *testing.T) {
	t.Parallel()

	// 既存の day1-am を更新する
	startAt := time.Date(2025, 6, 1, 9, 0, 0, 0, time.UTC)
	endAt := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)

	err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
		return tx.SaveSchedule(t.Context(), []*domain.ScheduleData{
			{Name: "day1-am", StartAt: startAt, EndAt: endAt},
		})
	})
	if err != nil {
		t.Fatalf("failed to save schedule: %+v", err)
	}

	actual, err := repo.GetSchedule(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// SaveSchedule deletes schedules not in the input, so only 1 should remain
	if len(actual) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(actual))
	}

	expected := &domain.ScheduleData{Name: "day1-am", StartAt: startAt, EndAt: endAt}
	if diff := cmp.Diff(expected, actual[0]); diff != "" {
		t.Errorf("schedule mismatch (-want +got):\n%s", diff)
	}
}
