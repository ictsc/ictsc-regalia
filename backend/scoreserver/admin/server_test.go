package admin_test

import (
	"os"
	"testing"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
)

func TestMain(m *testing.M) {
	run := pgtest.WrapRun(m.Run)
	os.Exit(run())
}

func assertCode(t *testing.T, expected connect.Code, got error) {
	t.Helper()

	var code connect.Code
	if got != nil {
		code = connect.CodeOf(got)
	}

	if code != expected {
		t.Errorf("expect code: %v, but got: %v", expected, got)
	}
}
