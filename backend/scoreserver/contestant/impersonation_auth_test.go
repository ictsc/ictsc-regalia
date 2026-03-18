package contestant

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
	adminauth "github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
)

func TestWithImpersonationAuthRejectsUnauthenticatedImpersonation(t *testing.T) {
	t.Parallel()

	called := false
	server, client := newImpersonationAuthTestServer(t, testAuthenticator{err: adminauth.ErrUnauthenticated}, func(r *http.Request) *session.UserSession {
		return &session.UserSession{
			UserID: uuid.Must(uuid.NewV4()),
			Impersonation: &session.ImpersonationSession{
				AdminName: "tester",
			},
		}
	}, &called)

	resp := doImpersonationAuthTestRequest(t, client, server.URL+"/protected")
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("resp.StatusCode = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("protected handler was called")
	}
}

func TestWithImpersonationAuthAllowsAuthenticatedImpersonation(t *testing.T) {
	t.Parallel()

	called := false
	server, client := newImpersonationAuthTestServer(t, testAuthenticator{viewer: adminauth.Viewer{
		Name:   "tester",
		Groups: []string{"tester"},
	}}, func(r *http.Request) *session.UserSession {
		return &session.UserSession{
			UserID: uuid.Must(uuid.NewV4()),
			Impersonation: &session.ImpersonationSession{
				AdminName: "tester",
			},
		}
	}, &called)

	resp := doImpersonationAuthTestRequest(t, client, server.URL+"/protected")
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("resp.StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
	if !called {
		t.Fatal("protected handler was not called")
	}
}

func TestWithImpersonationAuthAllowsNormalContestantSession(t *testing.T) {
	t.Parallel()

	called := false
	server, client := newImpersonationAuthTestServer(t, testAuthenticator{err: adminauth.ErrUnauthenticated}, func(r *http.Request) *session.UserSession {
		return &session.UserSession{
			UserID: uuid.Must(uuid.NewV4()),
		}
	}, &called)

	resp := doImpersonationAuthTestRequest(t, client, server.URL+"/protected")
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("resp.StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
	if !called {
		t.Fatal("protected handler was not called")
	}
}

func TestWithImpersonationAuthRejectsAuthenticatedViewerWithoutPermission(t *testing.T) {
	t.Parallel()

	called := false
	server, client := newImpersonationAuthTestServerWithPolicy(
		t,
		"g, tester, role:reader\np, role:reader, contestants, list",
		testAuthenticator{viewer: adminauth.Viewer{
			Name:   "tester",
			Groups: []string{"tester"},
		}},
		func(r *http.Request) *session.UserSession {
			return &session.UserSession{
				UserID: uuid.Must(uuid.NewV4()),
				Impersonation: &session.ImpersonationSession{
					AdminName: "tester",
				},
			}
		},
		&called,
	)

	resp := doImpersonationAuthTestRequest(t, client, server.URL+"/protected")
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("resp.StatusCode = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
	if called {
		t.Fatal("protected handler was called")
	}
}

type testAuthenticator struct {
	viewer adminauth.Viewer
	err    error
}

func (a testAuthenticator) HandleRequest(*http.Request) (*adminauth.Viewer, error) {
	if a.err != nil {
		return nil, a.err
	}
	viewer := a.viewer
	return &viewer, nil
}

func newImpersonationAuthTestServer(
	t *testing.T,
	authenticator testAuthenticator,
	sessionFactory func(*http.Request) *session.UserSession,
	called *bool,
) (*httptest.Server, *http.Client) {
	t.Helper()
	return newImpersonationAuthTestServerWithPolicy(t, "g, tester, role:admin", authenticator, sessionFactory, called)
}

func newImpersonationAuthTestServerWithPolicy(
	t *testing.T,
	policy string,
	authenticator testAuthenticator,
	sessionFactory func(*http.Request) *session.UserSession,
	called *bool,
) (*httptest.Server, *http.Client) {
	t.Helper()

	enforcer, err := adminauth.NewEnforcer(config.AdminAuthz{Policy: policy})
	if err != nil {
		t.Fatalf("failed to create enforcer: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/set", func(w http.ResponseWriter, r *http.Request) {
		if err := session.UserSessionStore.Write(r, w, sessionFactory(r), userSessionOption(nil, r)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusNoContent)
	})

	handler := http.Handler(mux)
	handler = withImpersonationAuth(nil, enforcer)(handler)
	handler = session.NewHandler(sessions.NewCookieStore([]byte("test-secret")))(handler)
	handler = adminauth.WithAuthn(handler, authenticator)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to create cookie jar: %v", err)
	}
	client := &http.Client{Jar: jar}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/auth/set", nil)
	if err != nil {
		t.Fatalf("failed to create session request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to initialize session: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		closeImpersonationAuthTestResponseBody(t, resp)
		t.Fatalf("failed to initialize session: status = %d", resp.StatusCode)
	}
	closeImpersonationAuthTestResponseBody(t, resp)

	return server, client
}

func doImpersonationAuthTestRequest(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	t.Cleanup(func() { closeImpersonationAuthTestResponseBody(t, resp) })
	return resp
}

func closeImpersonationAuthTestResponseBody(t *testing.T, resp *http.Response) {
	t.Helper()

	if err := resp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}
