package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
)

func Test_PgRepo_InvitationCode(t *testing.T) {
	t.Parallel()

	//nolint:thelper //ここではテストケースを書いているため
	tests := map[string]func(t *testing.T, db *sqlx.DB){
		"create": func(t *testing.T, db *sqlx.DB) {
			ctx, cancel := context.WithCancel(t.Context())
			t.Cleanup(cancel)

			repo := pg.NewRepository(db)

			team2, err := repo.GetTeamByCode(ctx, 2)
			if err != nil {
				t.Fatalf("failed to get team2: %+v", err)
			}

			code := &domain.InvitationCodeData{
				ID:        uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
				Code:      "ABCD1234EFGH5678",
				Team:      team2,
				ExpiresAt: must(time.Parse(time.RFC3339, "2025-01-02T00:00:00Z")),
				CreatedAt: must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")),
			}

			if err := repo.CreateInvitationCode(ctx, code); err != nil {
				t.Fatalf("failed to create invitation code: %+v", err)
			}

			row := db.QueryRowxContext(ctx, "SELECT * FROM invitation_codes WHERE id = $1", code.ID)
			if err := row.Err(); err != nil {
				t.Fatalf("failed to get invitation code: %+v", err)
			}
			got := map[string]any{}
			if err := row.MapScan(got); err != nil {
				t.Fatalf("failed to scan invitation code: %+v", err)
			}
			if diff := cmp.Diff(got, map[string]any{
				"id":         "00000000-0000-0000-0000-000000000001",
				"code":       "ABCD1234EFGH5678",
				"team_id":    "83027d5e-fa32-41d6-b290-fc38ba337f89",
				"expires_at": must(time.Parse(time.RFC3339, "2025-01-02T00:00:00Z")),
				"created_at": must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z")),
			}); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
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

func TestListInvitationCodes(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))
	got, err := repo.ListInvitationCodes(t.Context(), domain.InvitationCodeFilter{})
	if err != nil {
		t.Fatalf("failed to list invitation codes: %+v", err)
	}
	snaptest.Match(t, got)
}

func TestGetInvitationCode(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		code string

		wantErr error
	}{
		"ok": {
			code: "LHNZXGSF7L59WCG9",
		},
		"not found": {
			code:    "notfound",
			wantErr: domain.ErrInvitationCodeNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			got, err := repo.GetInvitationCode(t.Context(), tt.code)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			snaptest.Match(t, got)
		})
	}
}
