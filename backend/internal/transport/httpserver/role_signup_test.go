package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

func TestDiscordRoleSignup(t *testing.T) {
	for _, tc := range []struct {
		name, guild string
		roles       []string
		want        int
	}{
		{"mapped", "guild", []string{"team-a"}, 204},
		{"missing", "guild", []string{"staff"}, 403},
		{"ambiguous", "guild", []string{"team-a", "team-b"}, 403},
		{"wrong-guild", "other", []string{"team-a"}, 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newContractFixture(t)
			f.service.Config.ContestantGuildID = "guild"
			f.service.Config.DiscordRoleTeams = map[string]int64{"team-a": 2, "team-b": 3}
			_, err := f.store.CreateTeam(context.Background(), core.Team{Code: 2, Name: "Team A", Organization: "ICTSC", MemberLimit: 1})
			if err != nil {
				t.Fatal(err)
			}
			signup := f.createSession(t, session.Data{Kind: session.KindSignup, Discord: core.DiscordIdentity{ID: "123", Username: "alice", DisplayName: "Alice"}, GuildID: tc.guild, RoleIDs: tc.roles}, time.Hour)
			if tc.want == 204 {
				view := f.request(t, "GET", "/api/v1/viewer", "", signup)
				if view.Code != 200 || !strings.Contains(view.Body.String(), `"registration_team"`) {
					t.Fatalf("viewer: %d %s", view.Code, view.Body.String())
				}
			}
			r := httptest.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(`{"name":"alice","display_name":"Alice"}`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Origin", testOrigin)
			r.AddCookie(signup)
			w := httptest.NewRecorder()
			f.handler.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("signup: %d %s", w.Code, w.Body.String())
			}
			if tc.want == 204 {
				c, err := f.store.GetContestantByDiscord(context.Background(), "123")
				if err != nil || c.TeamCode != 2 {
					t.Fatalf("contestant: %+v %v", c, err)
				}
				_, err = f.store.RegisterContestant(context.Background(), core.Contestant{Name: "bob", DiscordID: "456", TeamCode: 2})
				if err == nil {
					t.Fatal("capacity not enforced")
				}
				w = httptest.NewRecorder()
				f.handler.ServeHTTP(w, r)
				if w.Code != 401 {
					t.Fatalf("used signup session: %d", w.Code)
				}
			}
		})
	}
}

func TestDiscordCallbackRejectsMissingTeamRole(t *testing.T) {
	f := newContractFixture(t)
	f.service.Config.AdminRoleIDs = map[string]struct{}{"other-role": {}}
	f.service.Config.ContestantGuildID = "guild-1"
	f.service.Config.DiscordRoleTeams = map[string]int64{"team-role": 2}
	token, err := f.sessions.Create(context.Background(), session.Data{Kind: session.KindOAuth, OAuthState: "test-state", PKCEVerifier: "verifier", Next: "/"}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/api/v1/auth/discord/callback?code=test&state=test-state", nil)
	r.AddCookie(&http.Cookie{Name: "oauth2-session", Value: token})
	w := httptest.NewRecorder()
	f.handler.ServeHTTP(w, r)
	if w.Code != 403 || !strings.Contains(w.Body.String(), "team_role_required") {
		t.Fatalf("callback: %d %s", w.Code, w.Body.String())
	}
}

func TestStaffContestantLoginRedirectsToAdmin(t *testing.T) {
	f := newContractFixture(t)
	f.service.Config.ContestantGuildID = "guild-1"
	f.service.Config.DiscordRoleTeams = map[string]int64{"team-role": 2}
	token, err := f.sessions.Create(context.Background(), session.Data{Kind: session.KindOAuth, OAuthState: "test-state", PKCEVerifier: "verifier", Next: "/"}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "/api/v1/auth/discord/callback?code=test&state=test-state", nil)
	r.AddCookie(&http.Cookie{Name: "oauth2-session", Value: token})
	w := httptest.NewRecorder()
	f.handler.ServeHTTP(w, r)
	if w.Code != 302 || w.Header().Get("Location") != "/admin/" {
		t.Fatalf("callback: %d %s", w.Code, w.Body.String())
	}
	var admin *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "admin-session" {
			admin = cookie
		}
	}
	if admin == nil {
		t.Fatal("missing admin cookie")
	}
	if _, err := f.sessions.Get(context.Background(), admin.Value, session.KindAdmin); err != nil {
		t.Fatal(err)
	}
}
