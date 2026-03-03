package pg_test

import (
	"context"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListAnswersForAdmin(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.ListAnswersForAdmin(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slices.SortStableFunc(actual, func(lhs, rhs *domain.AnswerData) int {
		return strings.Compare(lhs.ID.String(), rhs.ID.String())
	})
	snaptest.Match(t, actual)
}

func TestListAnswersByTeamProblem(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		viewer      string
		teamCode    int64
		problemCode string
	}{
		"ok/admin": {
			viewer:      "admin",
			teamCode:    1,
			problemCode: "AAA",
		},
		"ok/public": {
			viewer:      "public",
			teamCode:    1,
			problemCode: "AAA",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))

			var lister func(context.Context, int64, string) ([]*domain.AnswerData, error)
			switch tt.viewer {
			case "admin":
				lister = repo.ListAnswersByTeamProblemForAdmin
			case "public":
				lister = repo.ListAnswersByTeamProblemForPublic
			default:
				t.Fatalf("unexpected viewer: %s", tt.viewer)
			}

			actual, err := lister(t.Context(), tt.teamCode, tt.problemCode)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			snaptest.Match(t, actual)
		})
	}
}

func TestGetLatestAnswersByTeamProblemForPublic(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		teamID    uuid.UUID
		problemID uuid.UUID

		wantErr error
	}{
		"exists": {
			teamID:    uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			problemID: uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
		},
		"not found": {
			teamID:    uuid.FromStringOrNil("83027d5e-fa32-41d6-b290-fc38ba337f89"),
			problemID: uuid.FromStringOrNil("24f6aef0-5dcd-4032-825b-d1b19174a6f2"),

			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.GetLatestAnswerByTeamProblemForPublic(t.Context(), tt.teamID, tt.problemID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			snaptest.Match(t, actual)
		})
	}
}

func TestGetAnswerDetailForAdmin(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		teamCode     int64
		problemCode  string
		answerNumber uint32

		wantErr error
	}{
		"ok": {
			teamCode:     1,
			problemCode:  "AAA",
			answerNumber: 1,
		},
		"not found": {
			teamCode:     1,
			problemCode:  "AAA",
			answerNumber: 3,

			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.GetAnswerDetailForAdmin(t.Context(), tt.teamCode, tt.problemCode, tt.answerNumber)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			snaptest.Match(t, actual)
		})
	}
}

func TestCreateAnswer(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in *domain.AnswerDetailData

		queries map[string]string
	}{
		"ok": {
			in: &domain.AnswerDetailData{
				Answer: &domain.AnswerData{
					ID:     uuid.FromStringOrNil("1bb8bf23-95e1-438c-b30a-1778383190dc"),
					Number: 3,
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
						Category:     "Network",
						RedeployRule: domain.RedeployRuleUnredeployable,
					},
					Author: &domain.UserData{
						ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
						Name: "alice",
					},
					CreatedAt: time.Date(2025, 2, 3, 1, 0, 0, 0, time.UTC),
					Interval:  20 * time.Minute,
				},
				Body: &domain.AnswerBodyData{
					Descriptive: &domain.DescriptiveAnswerBodyData{
						Body: "answer",
					},
				},
			},
			queries: map[string]string{
				"answer": `
					SELECT 1 FROM answers WHERE
						id = '1bb8bf23-95e1-438c-b30a-1778383190dc' AND
						number = 3 AND
						team_id = 'a1de8fe6-26c8-42d7-b494-dea48e409091' AND
						problem_id = '16643c32-c686-44ba-996b-2fbe43b54513' AND
						user_id = '3a4ca027-5e02-4ade-8e2d-eddb39adc235' AND
						created_at_range = tstzrange('2025-02-03 01:00:00', '2025-02-03 01:20:00')`,
				"descriptive answer": `
					SELECT 1 FROM descriptive_answers WHERE
						answer_id = '1bb8bf23-95e1-438c-b30a-1778383190dc' AND
						body = 'answer'`,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)

			err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
				return tx.CreateAnswer(t.Context(), tt.in)
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for qn, query := range tt.queries {
				var dest any
				if err := db.GetContext(t.Context(), &dest, query); err != nil {
					t.Errorf("query %s: %v", qn, err)
				}
			}
		})
	}
}
