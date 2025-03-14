package pg_test

import (
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListTeamProblemScores(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		isPublic bool
	}{
		"admin":  {isPublic: false},
		"public": {isPublic: true},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))

			actual, err := repo.ListTeamProblemScores(t.Context(), tt.isPublic)
			if err != nil {
				t.Fatalf("ListTeamProblemScores failed: %v", err)
			}

			snaptest.Match(t, actual)
		})
	}
}

func TestListProblemScoresByTeamID(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		isPublic bool
		teamID   uuid.UUID
	}{
		"admin":  {isPublic: false, teamID: uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091")},
		"public": {isPublic: true, teamID: uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091")},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))

			actual, err := repo.ListTeamProblemScoresByTeamID(t.Context(), tt.isPublic, tt.teamID)
			if err != nil {
				t.Fatalf("ListProblemScoresByTeamID failed: %v", err)
			}

			snaptest.Match(t, actual)
		})
	}
}

func TestGetProblemScoreByTeamIDAndProblemID(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		isPublic  bool
		teamID    uuid.UUID
		problemID uuid.UUID

		wantErr error
	}{
		"admin": {
			isPublic:  false,
			teamID:    uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			problemID: uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
		},
		"public": {
			isPublic:  true,
			teamID:    uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			problemID: uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
		},
		"not found": {
			isPublic:  false,
			teamID:    uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
			problemID: uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000"),
			wantErr:   domain.ErrNotFound,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))

			actual, err := repo.GetTeamProblemScore(t.Context(), tt.isPublic, tt.teamID, tt.problemID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetProblemScoreByTeamIDAndProblemID failed: %v", err)
			}
			if err != nil {
				return
			}

			snaptest.Match(t, actual)
		})
	}
}
