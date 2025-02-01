package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
)

func Test_PgRepo_InvitationCode(t *testing.T) {
	t.Parallel()

	now := must(time.Parse(time.RFC3339, "2025-01-01T00:00:00Z"))

	team1 := must(domain.NewTeam(domain.TeamInput{
		ID:           must(uuid.NewV4()),
		Code:         1,
		Name:         "team1",
		Organization: "org1",
	}))
	invitationCode := must(domain.NewInvitationCode(domain.InvitationCodeInput{
		ID:        must(uuid.NewV4()),
		Code:      "ABCD1234EFGH5678",
		Team:      team1,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}))

	//nolint:thelper //ここではテストケースを書いているため
	tests := map[string]func(t *testing.T, db *sqlx.DB){
		"create": func(t *testing.T, db *sqlx.DB) {
			repo := pg.NewRepository(db)
			if err := repo.CreateInvitationCode(context.Background(), invitationCode); err != nil {
				t.Fatalf("failed to create invitation code: %+v", err)
			}

			row := db.QueryRowx("SELECT * FROM invitation_codes WHERE id = $1", invitationCode.ID())
			if row.Err() != nil {
				t.Fatalf("failed to get invitation code: %+v", row.Err())
			}
			got := map[string]any{}
			if err := row.MapScan(got); err != nil {
				t.Fatalf("failed to map scan: %+v", err)
			}
			if diff := cmp.Diff(got, map[string]any{
				"id":         invitationCode.ID().String(),
				"code":       invitationCode.Code(),
				"team_id":    invitationCode.Team().ID().String(),
				"expires_at": invitationCode.ExpiresAt(),
				"created_at": invitationCode.CreatedAt(),
			}); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		},
		"list": func(t *testing.T, db *sqlx.DB) {
			repo := pg.NewRepository(db)

			before, err := repo.ListInvitationCodes(context.Background(), domain.InvitationCodeFilter{})
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if len(before) != 0 {
				t.Errorf("unexpected invitation codes: %+v", before)
			}

			if err := repo.CreateInvitationCode(context.Background(), invitationCode); err != nil {
				t.Fatalf("failed to create invitation code: %+v", err)
			}

			got, err := repo.ListInvitationCodes(context.Background(), domain.InvitationCodeFilter{})
			if err != nil {
				t.Fatalf("failed to list invitation codes: %+v", err)
			}
			if diff := cmp.Diff(
				got, []*domain.InvitationCode{invitationCode},
				cmp.AllowUnexported(domain.InvitationCode{}, domain.Team{}),
			); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		},
	}

	setup := func(t *testing.T, db *sqlx.DB) {
		t.Helper()

		repo := pg.NewRepository(db)
		if err := repo.CreateTeam(context.Background(), team1); err != nil {
			t.Fatalf("failed to create team: %+v", err)
		}
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db, ok := pgtest.SetupDB(t)
			if !ok {
				t.FailNow()
			}
			setup(t, db)
			test(t, db)
		})
	}
}
