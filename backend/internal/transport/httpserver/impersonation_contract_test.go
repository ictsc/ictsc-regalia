package httpserver

import (
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAdminImpersonationSessionContract(t *testing.T) {
	f := newContractFixture(t)
	f.seedContestant(t)
	admin := f.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "staff"}, time.Hour)
	for _, authenticated := range []bool{false, true} {
		r := httptest.NewRequest("POST", "/api/v1/admin/impersonations", strings.NewReader(`{"contestant_name":"alice"}`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", testOrigin)
		if authenticated {
			r.AddCookie(admin)
		}
		w := httptest.NewRecorder()
		f.handler.ServeHTTP(w, r)
		if !authenticated {
			if w.Code != 401 {
				t.Fatalf("anonymous: %d %s", w.Code, w.Body.String())
			}
			continue
		}
		if w.Code != 204 {
			t.Fatalf("impersonation: %d %s", w.Code, w.Body.String())
		}
		var user *http.Cookie
		for _, cookie := range w.Result().Cookies() {
			if cookie.Name == "user-session" {
				user = cookie
			}
		}
		if user == nil || !user.Secure || !user.HttpOnly || user.Path != "/" || user.Domain != "" {
			t.Fatalf("invalid impersonation cookie flags")
		}
		viewer := f.request(t, "GET", "/api/v1/viewer", "", user)
		if viewer.Code != 200 || !strings.Contains(viewer.Body.String(), `"impersonated_by":"staff"`) || !strings.Contains(viewer.Body.String(), `"name":"alice"`) {
			t.Fatalf("viewer: %d %s", viewer.Code, viewer.Body.String())
		}
	}
}
