package pg_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
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

			actual, err := collectErrIter(repo.ListUsers(t.Context(), tt.filter))
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

func Test_PgRepo_GetDiscordLinkedUser(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		discordUserID int64

		want    *domain.UserData
		wantErr error
	}{
		"ok": {
			discordUserID: 123456789012345678,
			want: &domain.UserData{
				ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
				Name: "alice",
			},
		},
		"not found": {
			discordUserID: 999999999999999999,
			wantErr:       domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))

			got, err := repo.GetDiscordLinkedUser(t.Context(), tt.discordUserID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want error type %v, but got %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
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
			ctx := t.Context()

			profile := &domain.UserProfileData{
				User: &domain.UserData{
					ID:   uuid.FromStringOrNil("ab072031-3cf8-4795-9902-01e9e7fdf0bc"),
					Name: "charlie",
				},
				Profile: &domain.ProfileData{
					DisplayName: "Charlie",
				},
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
		"duplicate name": func(t *testing.T, repo *pg.Repository, _ *sqlx.DB) {
			ctx := t.Context()

			profile := &domain.UserProfileData{
				User: &domain.UserData{
					ID:   uuid.FromStringOrNil("1aacaafb-0f63-4746-be40-3b3511844c73"),
					Name: "alice",
				},
				Profile: &domain.ProfileData{
					DisplayName: "Alice",
				},
			}

			err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.CreateUser(ctx, profile)
			})

			if !errors.Is(err, domain.ErrDuplicateUserName) {
				t.Errorf("want error type %v, but got %v", domain.ErrDuplicateUserName, err)
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
func Test_PgRepo_LinkDiscordUser(t *testing.T) {
	t.Parallel()

	//nolint:thelper //サブテスト
	tests := map[string]func(t *testing.T, repo *pg.Repository, db *sqlx.DB){
		"ok": func(t *testing.T, repo *pg.Repository, db *sqlx.DB) {
			ctx := t.Context()
			userID := uuid.FromStringOrNil("c8c698ab-3b98-41a9-83b7-3cedd480112b")
			discordUserID := int64(123)

			if _, err := db.ExecContext(ctx, db.Rebind(`
				INSERT INTO users (id, name) VALUES (?, 'tarou')
			`), userID); err != nil {
				t.Fatalf("failed to insert user: %+v", err)
			}

			if err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.LinkDiscordUser(ctx, userID, discordUserID)
			}); err != nil {
				t.Fatalf("failed to link discord user: %+v", err)
			}

			row := db.QueryRowxContext(ctx, `
				SELECT user_id, discord_user_id
				FROM discord_users
				WHERE discord_user_id = $1
			`, discordUserID)
			if err := row.Err(); err != nil {
				t.Fatalf("failed to get discord user: %+v", err)
			}

			actual, expected := map[string]any{}, map[string]any{
				"user_id":         userID.String(),
				"discord_user_id": discordUserID,
			}
			if err := row.MapScan(actual); err != nil {
				t.Fatalf("failed to scan discord user: %+v", err)
			}
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		},
		"already exists": func(t *testing.T, repo *pg.Repository, _ *sqlx.DB) {
			ctx := t.Context()
			userID := uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235")
			discordUserID := int64(234567890123456789)

			err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.LinkDiscordUser(ctx, userID, discordUserID)
			})

			if !errors.Is(err, domain.ErrAlreadyExists) {
				t.Errorf("want error type %v, but got %v", domain.ErrAlreadyExists, err)
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

func TestGetUserProfileByID(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		userID uuid.UUID

		want    *domain.UserProfileData
		wantErr error
	}{
		"ok": {
			userID: uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
			want: &domain.UserProfileData{
				User: &domain.UserData{
					ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
					Name: "alice",
				},
				Profile: &domain.ProfileData{
					DisplayName: "Alice",
				},
			},
		},
		"not found": {
			userID:  uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo := pg.NewRepository(pgtest.SetupDB(t))
			got, err := repo.GetUserProfileByID(t.Context(), tt.userID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want error type %v, but got %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestListTeamMembers(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		want []*domain.TeamMemberProfileData
	}{
		"all": {
			want: []*domain.TeamMemberProfileData{
				{
					User: &domain.UserData{
						ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
						Name: "alice",
					},
					Team: &domain.TeamData{
						ID:           uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
						Code:         1,
						Name:         "トラブルシューターズ",
						Organization: "ICTSC Association",
						MaxMembers:   6,
					},
					Profile: &domain.ProfileData{
						DisplayName: "Alice",
					},
					DiscordUserID: 123456789012345678,
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.ListTeamMembers(t.Context())
			if err != nil {
				return
			}
			slices.SortStableFunc(actual, func(l, r *domain.TeamMemberProfileData) int {
				return strings.Compare(l.User.Name, r.User.Name)
			})
			if diff := cmp.Diff(tt.want, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestListTeamMembersByTeamID(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		teamID uuid.UUID
		want   []*domain.TeamMemberProfileData
	}{
		"ok": {
			teamID: uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			want: []*domain.TeamMemberProfileData{
				{
					User: &domain.UserData{
						ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
						Name: "alice",
					},
					Team: &domain.TeamData{
						ID:           uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
						Code:         1,
						Name:         "トラブルシューターズ",
						Organization: "ICTSC Association",
						MaxMembers:   6,
					},
					Profile: &domain.ProfileData{
						DisplayName: "Alice",
					},
					DiscordUserID: 123456789012345678,
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.ListTeamMembersByTeamID(t.Context(), tt.teamID)
			if err != nil {
				return
			}
			slices.SortStableFunc(actual, func(l, r *domain.TeamMemberProfileData) int {
				return strings.Compare(l.User.Name, r.User.Name)
			})
			if diff := cmp.Diff(tt.want, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
