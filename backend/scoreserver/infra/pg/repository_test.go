package pg_test

import (
	"os"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
)

func TestMain(m *testing.M) {
	run := pgtest.WrapRun(m.Run)
	os.Exit(run())
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
