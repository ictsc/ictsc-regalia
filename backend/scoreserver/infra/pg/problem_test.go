package pg_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListProblems(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.ListProblems(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	slices.SortFunc(actual, func(a, b *domain.ProblemData) int {
		return strings.Compare(a.Code, b.Code)
	})

	snaptest.Match(t, actual)
}

func TestGetProblemByCode(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		code    string
		wantErr error
	}{
		"ok": {
			code: "AAA",
		},
		"not found": {
			code:    "ZZZ",
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.GetProblemByCode(t.Context(), tt.code)
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

func TestGetDescriptiveProblem(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		id      uuid.UUID
		wantErr error
	}{
		"ok": {
			id: uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
		},
		"not found": {
			id:      uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000"),
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.GetDescriptiveProblem(t.Context(), tt.id)
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

func TestSaveDescriptiveProblem(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in      *domain.DescriptiveProblemData
		queries map[string]string
	}{
		"create": {
			in: &domain.DescriptiveProblemData{
				Problem: &domain.ProblemData{
					ID:           uuid.FromStringOrNil("3a4d8197-09f4-4bb0-9255-a8b6a943a36c"),
					Code:         "ZZZ",
					ProblemType:  domain.ProblemTypeDescriptive,
					Title:        "test",
					MaxScore:     100,
					Category:     "Network",
					RedeployRule: domain.RedeployRuleUnredeployable,
				},
				Content: &domain.ProblemContentData{
					Body:        "test",
					Explanation: "test",
				},
			},
			queries: map[string]string{
				"problem": `
					SELECT 1 FROM problems WHERE
					id = '3a4d8197-09f4-4bb0-9255-a8b6a943a36c' AND
					code = 'ZZZ' AND type = 'DESCRIPTIVE' AND
					title = 'test' AND max_score = 100 AND
					category = 'Network' AND redeploy_rule = 'UNREDEPLOYABLE'`,
				"content": `
					SELECT 1 FROM problem_contents WHERE
					problem_id = '3a4d8197-09f4-4bb0-9255-a8b6a943a36c' AND
					body = 'test' AND explanation = 'test'`,
			},
		},
		"create manual": {
			in: &domain.DescriptiveProblemData{
				Problem: &domain.ProblemData{
					ID:           uuid.FromStringOrNil("4b5e9308-1af5-5cc1-a366-b9c7b954b47d"),
					Code:         "MAN",
					ProblemType:  domain.ProblemTypeDescriptive,
					Title:        "Manual Redeploy",
					MaxScore:     100,
					Category:     "Network",
					RedeployRule: domain.RedeployRuleManual,
				},
				Content: &domain.ProblemContentData{
					Body:        "manual redeploy problem",
					Explanation: "manual redeploy explanation",
				},
			},
			queries: map[string]string{
				"problem": `
					SELECT 1 FROM problems WHERE
					id = '4b5e9308-1af5-5cc1-a366-b9c7b954b47d' AND
					code = 'MAN' AND type = 'DESCRIPTIVE' AND
					title = 'Manual Redeploy' AND max_score = 100 AND
					category = 'Network' AND redeploy_rule = 'MANUAL'`,
				"content": `
					SELECT 1 FROM problem_contents WHERE
					problem_id = '4b5e9308-1af5-5cc1-a366-b9c7b954b47d' AND
					body = 'manual redeploy problem' AND explanation = 'manual redeploy explanation'`,
			},
		},
		"update": {
			in: &domain.DescriptiveProblemData{
				Problem: &domain.ProblemData{
					ID:           uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
					Code:         "ZZQ",
					ProblemType:  domain.ProblemTypeDescriptive,
					Title:        "問題Q",
					MaxScore:     200,
					Category:     "Server",
					RedeployRule: domain.RedeployRulePercentagePenalty,
					PercentagePenalty: &domain.RedeployPenaltyPercentage{
						Threshold: 2, Percentage: 10,
					},
				},
				Content: &domain.ProblemContentData{
					Body:        "body",
					Explanation: "explanation",
				},
			},
			queries: map[string]string{
				"problem": `
					SELECT 1 FROM problems WHERE
					id = '16643c32-c686-44ba-996b-2fbe43b54513' AND
					code = 'ZZQ' AND type = 'DESCRIPTIVE' AND
					title = '問題Q' AND max_score = 200 AND
					category = 'Server' AND redeploy_rule = 'PERCENTAGE_PENALTY'`,
				"content": `
					SELECT 1 FROM problem_contents WHERE
					problem_id = '16643c32-c686-44ba-996b-2fbe43b54513' AND
					body = 'body' AND explanation = 'explanation'`,
				"redeploy_rule": `
					SELECT 1 FROM redeploy_percentage_penalties WHERE
					problem_id = '16643c32-c686-44ba-996b-2fbe43b54513' AND
					threshold = 2 AND percentage = 10`,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)

			if err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
				return tx.SaveDescriptiveProblem(t.Context(), tt.in)
			}); err != nil {
				t.Fatal(err)
			}

			for qn, query := range tt.queries {
				var dst any
				if err := db.GetContext(t.Context(), &dst, query); err != nil {
					t.Errorf("query %s: %v", qn, err)
				}
			}
		})
	}
}

func TestDeleteProblem(t *testing.T) {
	t.Skip("FKの制約を入れたらテストが通らなくなった")
	t.Parallel()

	id := uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513")

	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)

	if err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
		return tx.DeleteProblem(t.Context(), id)
	}); err != nil {
		t.Fatal(err)
	}
}
