package pg_test

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
)

func Test_PgRepo_ListUsers(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		filter domain.UserListFilter

		wants []*domain.UserData
	}{
		"all": {
			wants: []*domain.UserData{
				{ID: uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"), Name: "alice"},
				{ID: uuid.FromStringOrNil("c4530ce6-d990-4414-8389-feca26883115"), Name: "bob"},
			},
		},
		"filter by name": {
			filter: domain.UserListFilter{
				Name: "alice",
			},
			wants: []*domain.UserData{
				{ID: uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"), Name: "alice"},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))

			actual, err := collectErrIter(repo.ListUsers(context.Background(), tt.filter))
			if err != nil {
				return
			}
			slices.SortStableFunc(actual, func(l, r *domain.UserData) int {
				return strings.Compare(l.Name, r.Name)
			})
			if diff := cmp.Diff(tt.wants, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_PgRepo_CreateUser(t *testing.T) {
	t.Parallel()

	//nolint:thelper //サブテスト
	tests := map[string]func(t *testing.T, repo *pg.Repository, db *sqlx.DB){
		"ok": func(t *testing.T, repo *pg.Repository, db *sqlx.DB) {
			ctx := context.Background()

			profile := &domain.UserProfileData{
				User: &domain.UserData{
					ID:   uuid.FromStringOrNil("ab072031-3cf8-4795-9902-01e9e7fdf0bc"),
					Name: "charlie",
				},
				DisplayName: "Charlie",
			}

			if err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.CreateUser(ctx, profile)
			}); err != nil {
				t.Fatalf("failed to create user profile: %+v", err)
			}

			row := db.QueryRowxContext(ctx, `
				SELECT
					u.id AS "u.id", u.name AS "u.name",
					p.display_name AS "p.display_name"
				FROM users AS u
				LEFT JOIN user_profiles AS p ON p.user_id = u.id
				WHERE u.id = 'ab072031-3cf8-4795-9902-01e9e7fdf0bc'`,
			)
			if err := row.Err(); err != nil {
				t.Fatalf("failed to get user: %+v", err)
			}

			actual, expected := map[string]any{}, map[string]any{
				"u.id":           "ab072031-3cf8-4795-9902-01e9e7fdf0bc",
				"u.name":         "charlie",
				"p.display_name": "Charlie",
			}
			if err := row.MapScan(actual); err != nil {
				t.Fatalf("failed to scan user: %+v", err)
			}
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		},
		"重複していたらエラー": func(t *testing.T, repo *pg.Repository, _ *sqlx.DB) {
			ctx := context.Background()

			profile := &domain.UserProfileData{
				User: &domain.UserData{
					ID:   uuid.FromStringOrNil("1aacaafb-0f63-4746-be40-3b3511844c73"),
					Name: "alice",
				},
				DisplayName: "Alice",
			}

			err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.CreateUser(ctx, profile)
			})

			expectedType, actualType := domain.ErrTypeAlreadyExists, domain.ErrTypeFrom(err)
			if expectedType != actualType {
				t.Errorf("want error type %v, but got %v", expectedType, actualType)
			}
		},
	}

	for name, tf := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)
			tf(t, repo, db)
		})
	}
}
