package httpserver

import (
	"encoding/json"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTeamColorContract(t *testing.T) {
	f := newContractFixture(t)
	admin := f.createSession(t, session.Data{Kind: session.KindAdmin}, time.Hour)
	request := func(method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", testOrigin)
		if cookie != nil {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		f.handler.ServeHTTP(w, r)
		return w
	}
	create := request("POST", "/api/v1/admin/teams", `{"code":12,"name":"Color Team","organization":"ICTSC","member_limit":4}`, admin)
	if create.Code != 201 {
		t.Fatalf("create: %d %s", create.Code, create.Body.String())
	}
	assertColor := func(w *httptest.ResponseRecorder, want string) {
		t.Helper()
		var body struct {
			Team struct {
				Color string `json:"color"`
			} `json:"team"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Team.Color != want {
			t.Fatalf("color = %q, want %q", body.Team.Color, want)
		}
	}
	assertColor(create, "#A6E35F")
	update := request("PATCH", "/api/v1/admin/teams/12", `{"color":"#0083C3"}`, admin)
	if update.Code != 200 {
		t.Fatalf("update: %d %s", update.Code, update.Body.String())
	}
	assertColor(update, "#0083C3")
	preserved := request("PATCH", "/api/v1/admin/teams/12", `{"name":"Renamed"}`, admin)
	assertColor(preserved, "#0083C3")
	invalid := request("PATCH", "/api/v1/admin/teams/12", `{"color":"#FFFFFF"}`, admin)
	if invalid.Code != 422 {
		t.Fatalf("invalid color status = %d", invalid.Code)
	}
	denied := request("PATCH", "/api/v1/admin/teams/12", `{"color":"#FFE000"}`, nil)
	if denied.Code != 401 {
		t.Fatalf("anonymous mutation status = %d", denied.Code)
	}
	assertColor(request("GET", "/api/v1/admin/teams/12", "", admin), "#0083C3")
}
