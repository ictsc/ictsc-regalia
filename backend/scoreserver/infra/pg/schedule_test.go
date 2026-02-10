package pg_test

import (
	"slices"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	startAt := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC)

	expected := []*domain.ScheduleData{
		{
			ID:      uuid.FromStringOrNil("8a23ce3f-4506-48e9-bf68-7d2d90592bf1"),
			StartAt: startAt,
			EndAt:   endAt,
		},
		{
			ID:      uuid.FromStringOrNil("4e72d440-dfde-4923-801d-0fd5ee2c0730"),
			StartAt: endAt,
			EndAt:   endAt.Add(24 * time.Hour),
		},
	}

	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.GetSchedule(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	slices.SortFunc(expected, func(a, b *domain.ScheduleData) int {
		return a.StartAt.Compare(b.StartAt)
	})
	slices.SortFunc(actual, func(a, b *domain.ScheduleData) int {
		return a.StartAt.Compare(b.StartAt)
	})

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("schedules mismatch (-want +got):\n%s", diff)
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
					ID:      uuid.FromStringOrNil("11d41cc7-c7c8-45c4-990b-0e3bcab4e54d"),
					StartAt: startAt,
					EndAt:   endAt,
				},
			},
			queries: []string{
				`SELECT 1
				 FROM schedules
				 WHERE id = '11d41cc7-c7c8-45c4-990b-0e3bcab4e54d'
				   AND start_at = '2025-02-03 00:00:00'
				   AND end_at = '2025-02-04 00:00:00'`,
			},
		},
		"multiple schedules": {
			input: []*domain.ScheduleData{
				{
					ID:      uuid.FromStringOrNil("816c496b-6bcb-46c2-b4df-a8c537bae51b"),
					StartAt: startAt,
					EndAt:   endAt,
				},
				{
					ID:      uuid.FromStringOrNil("5b88ebcc-2ef6-4ae7-b401-18eefe72cbbe"),
					StartAt: endAt,
					EndAt:   endAt.Add(24 * time.Hour),
				},
			},
			queries: []string{
				`SELECT 1
				FROM schedules
				WHERE id = '816c496b-6bcb-46c2-b4df-a8c537bae51b'
				AND start_at = '2025-02-03 00:00:00'
				AND end_at = '2025-02-04 00:00:00'`,
				`SELECT 1
				FROM schedules
				WHERE id = '5b88ebcc-2ef6-4ae7-b401-18eefe72cbbe'
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
