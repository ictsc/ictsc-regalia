package contestant_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
)

func TestSignUpOK(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	input := contestant.SignUpInput{
		InvitationCode: "LHNZXGSF7L59WCG9",
		Name:           "test",
		DisplayName:    "Tester",
		DiscordID:      "987654321987654321",
	}

	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)
	effect := pg.Tx(repo, func(rt *pg.RepositoryTx) contestant.SignUpTxEffect { return rt })
	userProfile, err := contestant.SignUp(t.Context(), effect, now, &input)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(
		userProfile.Data(), &domain.UserProfileData{
			User: &domain.UserData{
				Name: "test",
			},
			DisplayName: "Tester",
		},
		cmpopts.IgnoreFields(domain.UserData{}, "ID"),
	); diff != "" {
		t.Errorf("diff: %s", diff)
	}

	var result int
	if err := sqlx.GetContext(t.Context(), db, &result, `
		SELECT 1 FROM team_members AS tm
		JOIN invitation_codes AS ic ON tm.invitation_code_id = ic.id
		WHERE tm.user_id = $1 AND ic.code = $2`,
		userProfile.User().ID(), "LHNZXGSF7L59WCG9",
	); err != nil {
		t.Fatal(err)
	}
	if err := sqlx.GetContext(t.Context(), db, &result, `
		SELECT 1 FROM discord_users WHERE user_id = $1 AND discord_user_id = $2`,
		userProfile.User().ID(), "987654321987654321",
	); err != nil {
		t.Fatal(err)
	}
}
