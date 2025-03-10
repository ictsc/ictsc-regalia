package pg_test

import (
	"slices"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListMarkingResults(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))

	actual, err := repo.ListMarkingResults(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	slices.SortStableFunc(actual, func(i, j *domain.MarkingResultData) int {
		return i.CreatedAt.Compare(j.CreatedAt)
	})
	snaptest.Match(t, actual)
}

func TestCreateMarkingResult(t *testing.T) {
	t.Parallel()

	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)

	var markCount int
	if err := db.Get(&markCount, "SELECT COUNT(1) FROM marking_results"); err != nil {
		t.Fatal(err)
	}

	if err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
		return tx.CreateMarkingResult(t.Context(), &domain.MarkingResultData{
			ID:    uuid.FromStringOrNil("53d6a2aa-ebf8-4bd6-98f2-baa9758ae6a6"),
			Judge: "judge",
			Answer: &domain.AnswerData{
				ID: uuid.FromStringOrNil("7cedf13e-5325-425e-a5d6-fea5fc127e49"),
			},
			Score: &domain.ScoreData{
				MarkedScore: 100,
			},
			Rationale: &domain.MarkingRationaleData{
				DescriptiveComment: "comment",
			},
			CreatedAt: time.Date(2025, 2, 3, 2, 0, 0, 0, time.UTC),
		})
	}); err != nil {
		t.Fatal(err)
	}

	var count int
	if err := db.Get(&count, `SELECT COUNT(*) FROM marking_results`); err != nil {
		t.Fatal(err)
	}
	if expected := markCount + 1; count != expected {
		t.Errorf("unexpected marking_result count: %d, expected %d", count, expected)
	}

	if err := db.Get(&count, `
		SELECT COUNT(*) FROM marking_results WHERE
			id = '53d6a2aa-ebf8-4bd6-98f2-baa9758ae6a6' AND
			answer_id = '7cedf13e-5325-425e-a5d6-fea5fc127e49' AND
			judge_name = 'judge' AND
			created_at = '2025-02-03 02:00:00+00'`,
	); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Error("marking_result not found")
	}
}
