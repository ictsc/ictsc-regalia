package pg_test

import (
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
)

func TestGetTeamMember(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inUserID uuid.UUID
		wantErr  error
	}{
		"ok": {
			inUserID: uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
		},
		"not found": {
			inUserID:/* bob */ uuid.FromStringOrNil("c4530ce6-d990-4414-8389-feca26883115"),
			wantErr: domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)

			got, err := repo.GetTeamMemberByID(ctx, tt.inUserID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr: %v, got: %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}

			snaptest.Match(t, got)
		})
	}
}

func TestCountTeamMembers(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inTeamID uuid.UUID
		wants    uint
		wantErr  error
	}{
		"ok": {
			inTeamID: uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			wants:    1,
		},
		"not found": {
			inTeamID: uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
			wants:    0,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)

			got, err := repo.CountTeamMembers(ctx, tt.inTeamID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wantErr: %v, got: %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}

			if got != tt.wants {
				t.Errorf("want: %v, got: %v", tt.wants, got)
			}
		})
	}
}

func TestAddTeamMember(t *testing.T) {
	t.Parallel()

	//nolint:thelper //テストケース
	tests := map[string]func(t *testing.T, db *sqlx.DB){
		"ok": func(t *testing.T, db *sqlx.DB) {
			// bob はどのチームにも所属していない
			userID := uuid.FromStringOrNil("c4530ce6-d990-4414-8389-feca26883115")
			teamID := uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091")
			codeID := uuid.FromStringOrNil("ad3f83d3-65be-4884-8a03-adb11a8127ef")

			ctx := t.Context()
			repo := pg.NewRepository(db)
			if err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.AddTeamMember(ctx, userID, codeID, teamID)
			}); err != nil {
				t.Fatalf("failed to add member: %v", err)
			}

			row := db.QueryRowxContext(ctx, "SELECT user_id, team_id, invitation_code_id FROM team_members WHERE user_id = $1", userID)
			if err := row.Err(); err != nil {
				t.Fatalf("failed to select: %v", err)
			}

			expected, actual := map[string]any{
				"user_id":            userID.String(),
				"team_id":            teamID.String(),
				"invitation_code_id": codeID.String(),
			}, map[string]any{}
			if err := row.MapScan(actual); err != nil {
				t.Fatalf("failed to map scan: %v", err)
			}

			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		},
		"already exists": func(t *testing.T, db *sqlx.DB) {
			// alice は既にトラブルシューターズに所属している
			userID := uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235")
			teamID := uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091")
			codeID := uuid.FromStringOrNil("ad3f83d3-65be-4884-8a03-adb11a8127ef")

			ctx := t.Context()
			repo := pg.NewRepository(db)

			err := repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
				return tx.AddTeamMember(ctx, userID, codeID, teamID)
			})
			if !errors.Is(err, domain.ErrAlreadyExists) {
				t.Errorf("wantErr: %v, got: %v", domain.ErrAlreadyExists, err)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db := pgtest.SetupDB(t)
			test(t, db)
		})
	}
}
