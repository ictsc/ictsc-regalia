package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
)

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

func setupEnforcer(tb testing.TB) *auth.Enforcer {
	tb.Helper()

	enforcer, err := auth.NewEnforcer(config.AdminAuthz{
		Policy: "g, tester, role:admin",
	})
	if err != nil {
		tb.Fatalf("Failed to create enforcer: %v", err)
	}
	return enforcer
}

func setupServer(tb testing.TB, handler http.Handler) *httptest.Server {
	tb.Helper()

	handler = auth.WithAuthn(handler, &staticAuthenticator{
		Viewer: auth.Viewer{
			Name:   "tester",
			Groups: []string{"tester"},
		},
	})
	server := httptest.NewServer(handler)
	tb.Cleanup(server.Close)
	return server
}

type staticAuthenticator struct {
	Viewer auth.Viewer
}

func (a *staticAuthenticator) HandleRequest(*http.Request) (*auth.Viewer, error) {
	v := a.Viewer
	return &v, nil
}
