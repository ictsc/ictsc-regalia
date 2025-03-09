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

	expected := []*domain.MarkingResultData{
		{
			ID:    uuid.FromStringOrNil("862b646a-5fdd-4a77-bb2d-7ef5d4f1d069"),
			Judge: "judge",
			Answer: &domain.AnswerData{
				ID:     uuid.FromStringOrNil("7cedf13e-5325-425e-a5d6-fea5fc127e49"),
				Number: 1,
				Team: &domain.TeamData{
					ID:           uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
					Code:         1,
					Name:         "トラブルシューターズ",
					Organization: "ICTSC Association",
					MaxMembers:   6,
				},
				Problem: &domain.ProblemData{
					ID:           uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
					Code:         "ZZA",
					ProblemType:  domain.ProblemTypeDescriptive,
					Title:        "問題A",
					MaxScore:     100,
					RedeployRule: domain.RedeployRuleUnredeployable,
				},
				Author: &domain.UserData{
					ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
					Name: "alice",
				},
				CreatedAt: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
				Interval:  20 * time.Minute,
			},
			Score: &domain.ScoreData{
				MarkedScore: 80,
			},
			Rationale: &domain.MarkingRationaleData{
				DescriptiveComment: "問題Aへのチーム1の解答1の採点理由",
			},
			CreatedAt: time.Date(2025, 2, 3, 1, 0, 0, 0, time.UTC),
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("(-expected, +actual)\n%s", diff)
	}
}

func TestCreateMarkingResult(t *testing.T) {
	t.Parallel()

	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)

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
	if count != 2 {
		t.Errorf("unexpected marking_result count: %d", count)
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
