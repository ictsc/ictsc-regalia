package pg_test

import (
	"context"
	"os"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(pgtest.WrapRun(m.Run)())
}

func Test_Team(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db, ok := pgtest.SetupDB(t)
	if !ok {
		return
	}

	//nolint:exhaustruct
	team := pg.Team{
		Code:         001,
		Name:         "Trouble Shooters",
		Organization: "ICTSC Committee",
	}

	if _, err := db.NamedExecContext(ctx, "INSERT INTO teams (code, name, organization) VALUES (:code, :name, :organization)", team); err != nil {
		require.NoError(t, err)
		return
	}

	var dbTeam pg.Team
	if err := db.GetContext(ctx, &dbTeam, "SELECT * FROM teams WHERE code = $1", team.Code); err != nil {
		require.NoError(t, err)
		return
	}

	assert.Equal(t, team.Code, dbTeam.Code)
	assert.Equal(t, team.Name, dbTeam.Name)
	assert.Equal(t, team.Organization, dbTeam.Organization)
}
