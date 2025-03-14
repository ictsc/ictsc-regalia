package pg_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListDeployments(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))

	deployments, err := repo.ListDeployments(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	slices.SortStableFunc(deployments, func(i, j *domain.DeploymentData) int {
		return strings.Compare(i.ID.String(), j.ID.String())
	})

	snaptest.Match(t, deployments)
}

func TestCreateDeployment(t *testing.T) {
	t.Parallel()

	act := func(
		t *testing.T, repo *pg.Repository, input *domain.CreateDeploymentInput,
	) error {
		t.Helper()

		return repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
			return tx.CreateDeployment(t.Context(), input)
		})
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		db := pgtest.SetupDB(t)
		repo := pg.NewRepository(db)

		var beforeCount int
		if err := db.GetContext(
			t.Context(), &beforeCount, "SELECT COUNT(*) FROM redeployment_requests",
		); err != nil {
			t.Fatal(err)
		}

		if err := act(t, repo, &domain.CreateDeploymentInput{
			ID:         uuid.FromStringOrNil("8c94f006-a2a7-451b-9c56-d81b70aabeaf"),
			TeamID:     uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			ProblemID:  uuid.FromStringOrNil("24f6aef0-5dcd-4032-825b-d1b19174a6f2"),
			Revision:   2,
			Status:     domain.DeploymentStatusQueued,
			OccurredAt: time.Date(2021, 1, 2, 3, 4, 5, 0, time.UTC),
		}); err != nil {
			t.Fatal(err)
		}

		var afterCount int
		if err := db.GetContext(
			t.Context(), &afterCount, "SELECT COUNT(*) FROM redeployment_requests",
		); err != nil {
			t.Fatal(err)
		}

		if afterCount-beforeCount != 1 {
			t.Errorf("unexpected count: got %d, want %d", afterCount, beforeCount+1)
		}

		var count int
		if err := db.GetContext(
			t.Context(), &count, `
			SELECT 1 FROM redeployment_requests
			WHERE id = '8c94f006-a2a7-451b-9c56-d81b70aabeaf' AND
				team_id = 'a1de8fe6-26c8-42d7-b494-dea48e409091' AND
				problem_id = '24f6aef0-5dcd-4032-825b-d1b19174a6f2' AND
				revision = 2`,
		); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Error("deployment not found")
		}
	})
}
