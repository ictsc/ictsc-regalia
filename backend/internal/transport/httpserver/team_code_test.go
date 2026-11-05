package httpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

func TestAdminTeamCodesRunFromOneToNinetyNine(t *testing.T) {
	f := newContractFixture(t)
	admin := f.createSession(t, session.Data{Kind: session.KindAdmin}, time.Hour)
	create := func(code string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/admin/teams",
			strings.NewReader(`{"code":`+code+`,"name":"Team `+code+`","organization":"ICTSC","member_limit":4}`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", testOrigin)
		r.AddCookie(admin)
		w := httptest.NewRecorder()
		f.handler.ServeHTTP(w, r)
		return w
	}
	if w := create("1"); w.Code != http.StatusCreated {
		t.Fatalf("team 1: %d %s", w.Code, w.Body.String())
	}
	for _, code := range []string{"0", "100"} {
		if w := create(code); w.Code != http.StatusUnprocessableEntity && w.Code != http.StatusBadRequest {
			t.Fatalf("team %s: %d %s, want a validation error", code, w.Code, w.Body.String())
		}
	}
}
