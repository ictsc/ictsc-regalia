package pg_test

import (
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestGetRanking(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		isPublic bool
	}{
		"public": {isPublic: true},
		"admin":  {isPublic: false},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repo := pg.NewRepository(pgtest.SetupDB(t))
			actual, err := repo.GetRanking(t.Context(), tt.isPublic)
			if err != nil {
				t.Fatal(err)
			}

			snaptest.Match(t, actual)
		})
	}
}
