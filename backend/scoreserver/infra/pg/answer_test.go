package pg_test

import (
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListAnswers(t *testing.T) {
	t.Parallel()

	expected := []*domain.AnswerData{
		{
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
		{
			ID:     uuid.FromStringOrNil("abbe9c4e-eef5-40ac-a04e-6d8877b15185"),
			Number: 1,
			Team: &domain.TeamData{
				ID:           uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
				Code:         1,
				Name:         "トラブルシューターズ",
				Organization: "ICTSC Association",
				MaxMembers:   6,
			},
			Problem: &domain.ProblemData{
				ID:           uuid.FromStringOrNil("24f6aef0-5dcd-4032-825b-d1b19174a6f2"),
				Code:         "ZZB",
				ProblemType:  domain.ProblemTypeDescriptive,
				Title:        "問題B",
				MaxScore:     200,
				RedeployRule: domain.RedeployRulePercentagePenalty,
				PercentagePenalty: &domain.RedeployPenaltyPercentage{
					Threshold:  2,
					Percentage: 10,
				},
			},
			Author: &domain.UserData{
				ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
				Name: "alice",
			},
			CreatedAt: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			Interval:  20 * time.Minute,
		},
		{
			ID:     uuid.FromStringOrNil("4bb7a232-e0de-4b6d-b1a3-8e50737d73b2"),
			Number: 2,
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
			CreatedAt: time.Date(2025, 2, 3, 0, 30, 0, 0, time.UTC),
			Interval:  20 * time.Minute,
		},
	}
	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.ListAnswers(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("unexpected result (-want +got):\n%s", diff)
	}
}

func TestListAnswersByTeamProblem(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		teamCode    int64
		problemCode string

		wants []*domain.AnswerData
	}{
		"ok": {
			teamCode:    1,
			problemCode: "ZZA",
			wants: []*domain.AnswerData{
				{
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
				{
					ID:     uuid.FromStringOrNil("4bb7a232-e0de-4b6d-b1a3-8e50737d73b2"),
					Number: 2,
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
					CreatedAt: time.Date(2025, 2, 3, 0, 30, 0, 0, time.UTC),
					Interval:  20 * time.Minute,
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.ListAnswersByTeamProblem(t.Context(), tt.teamCode, tt.problemCode)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.wants, actual); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetLatestAnswersByTeamProblem(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		teamID    uuid.UUID
		problemID uuid.UUID

		wantErr error
		wants   *domain.AnswerData
	}{
		"exists": {
			teamID:    uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			problemID: uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),

			wants: &domain.AnswerData{
				ID:     uuid.FromStringOrNil("4bb7a232-e0de-4b6d-b1a3-8e50737d73b2"),
				Number: 2,
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
				CreatedAt: time.Date(2025, 2, 3, 0, 30, 0, 0, time.UTC),
				Interval:  20 * time.Minute,
			},
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
			actual, err := repo.GetLatestAnswerByTeamProblem(t.Context(), tt.teamID, tt.problemID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.wants, actual); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetAnswerDetail(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		teamCode     int64
		problemCode  string
		answerNumber uint32

		wantErr error
		wants   *domain.AnswerDetailData
	}{
		"ok": {
			teamCode:     1,
			problemCode:  "ZZA",
			answerNumber: 1,

			wants: &domain.AnswerDetailData{
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
				Body: &domain.AnswerBodyData{
					Descriptive: &domain.DescriptiveAnswerBodyData{
						Body: "問題Aへのチーム1の回答1",
					},
				},
			},
		},
		"not found": {
			teamCode:     1,
			problemCode:  "ZZA",
			answerNumber: 3,

			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.GetAnswerDetail(t.Context(), tt.teamCode, tt.problemCode, tt.answerNumber)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.wants, actual); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
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
						created_at = '2025-02-03 01:00:00' AND
						rate_limit_interval = '20m' AND
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
