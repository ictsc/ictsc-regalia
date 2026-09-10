package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

const (
	testOrigin         = "https://score.example.test"
	testSStateToken    = "fixed-sstate-callback-token"
	testContentCommit  = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testContestantName = "alice"
)

var testNow = time.Date(2026, 8, 30, 6, 0, 0, 0, time.UTC)

type contractFixture struct {
	handler  http.Handler
	service  *service.Service
	store    *memory.CompetitionStore
	sessions *memory.SessionStore
}

type fakeDiscord struct{}

func (fakeDiscord) AuthorizationURL(state, codeChallenge, redirectURI string, admin bool) string {
	query := url.Values{
		"state":          {state},
		"code_challenge": {codeChallenge},
		"redirect_uri":   {redirectURI},
	}
	if admin {
		query.Set("admin", "true")
	}
	return "https://discord.example.test/oauth2/authorize?" + query.Encode()
}

func (fakeDiscord) Exchange(context.Context, string, string, string, bool) (service.DiscordResult, error) {
	return service.DiscordResult{
		Identity: core.DiscordIdentity{ID: "discord-1", Username: "alice", DisplayName: "Alice"},
		GuildID:  "guild-1",
		RoleIDs:  []string{"admin-role"},
	}, nil
}

func newContractFixture(t *testing.T) *contractFixture {
	t.Helper()
	store := memory.NewCompetitionStore()
	sessions := memory.NewSessionStore()
	svc := service.New(store, sessions, service.Config{
		ContestantRedirectURI: "https://score.example.test/api/v1/auth/discord/callback",
		AdminRedirectURI:      "https://score.example.test/api/v1/admin/auth/discord/callback",
		AdminGuildID:          "guild-1",
		AdminRoleIDs:          map[string]struct{}{"admin-role": {}},
	})
	svc.Discord = fakeDiscord{}
	svc.Events = memory.NewEventBus()
	svc.Now = func() time.Time { return testNow }
	handler, err := New(svc, Options{
		SecureCookies:       true,
		AllowedOrigins:      []string{testOrigin},
		SStateCallbackToken: testSStateToken,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return &contractFixture{handler: handler, service: svc, store: store, sessions: sessions}
}

func (f *contractFixture) request(t *testing.T, method, target, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	recorder := httptest.NewRecorder()
	f.handler.ServeHTTP(recorder, request)
	return recorder
}

func (f *contractFixture) createSession(t *testing.T, data session.Data, ttl time.Duration) *http.Cookie {
	t.Helper()
	token, err := f.sessions.Create(context.Background(), data, ttl)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	name := map[session.Kind]string{
		session.KindContestant: "regalia-user-session",
		session.KindAdmin:      "admin-session",
		session.KindSignup:     "signup-session",
	}[data.Kind]
	return &http.Cookie{Name: name, Value: token}
}

func (f *contractFixture) seedContestant(t *testing.T) core.Contestant {
	t.Helper()
	ctx := context.Background()
	team := core.Team{Code: 2, Name: "Team Two", Organization: "Example University", MemberLimit: 4}
	if _, err := f.store.CreateTeam(ctx, team); err != nil {
		t.Fatalf("create team: %v", err)
	}
	invitation := core.Invitation{Code: "invite-alice", TeamCode: team.Code, CreatedAt: testNow.Add(-time.Hour), ExpiresAt: testNow.Add(time.Hour)}
	if _, err := f.store.CreateInvitation(ctx, invitation); err != nil {
		t.Fatalf("create invitation: %v", err)
	}
	contestant, err := f.store.ConsumeInvitation(ctx, invitation.Code, testNow, core.Contestant{
		Name: testContestantName, DisplayName: "Alice", DiscordID: "discord-1",
	})
	if err != nil {
		t.Fatalf("consume invitation: %v", err)
	}
	return contestant
}

func (f *contractFixture) seedActiveContent(t *testing.T) {
	t.Helper()
	zero := int32(0)
	snapshot := core.ContentSnapshot{
		CommitSHA:  testContentCommit,
		Repository: "ictsc/example-content",
		Ref:        "refs/heads/main",
		FetchedAt:  testNow,
		Manifest: core.Manifest{
			Version: 1,
			Sections: []core.Section{{
				Slug: "day-1", Beginning: testNow.Add(-time.Hour), Ending: testNow.Add(time.Hour), ProblemIDs: []string{"P1"},
			}},
			Problems: []core.Problem{{
				Code: "P1", Title: "Problem One", MaxScore: 100, Category: "network", SectionSlug: "day-1",
				Type: "DESCRIPTIVE", Body: "body", Redeploy: core.RedeployRule{Type: core.RedeployPercentage, Threshold: &zero, Percentage: &zero},
			}},
		},
	}
	if _, err := f.store.ActivateContent(context.Background(), snapshot, ""); err != nil {
		t.Fatalf("activate content: %v", err)
	}
}

func TestAnonymousSignoutIsIdempotentAndClearsAllTransientCookies(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		cookies map[string]http.SameSite
	}{
		{name: "contestant", path: "/api/v1/auth/signout", cookies: map[string]http.SameSite{
			"oauth2-session": http.SameSiteLaxMode, "signup-session": http.SameSiteStrictMode, "regalia-user-session": http.SameSiteStrictMode,
		}},
		{name: "admin", path: "/api/v1/admin/auth/signout", cookies: map[string]http.SameSite{
			"admin-oauth2-session": http.SameSiteLaxMode, "admin-session": http.SameSiteStrictMode,
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newContractFixture(t)
			request := httptest.NewRequest(http.MethodPost, test.path, nil)
			request.Header.Set("Origin", testOrigin)
			recorder := httptest.NewRecorder()
			fixture.handler.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want 204; body = %s", recorder.Code, recorder.Body.String())
			}

			response := recorder.Result()
			got := make(map[string]*http.Cookie)
			for _, cookie := range response.Cookies() {
				got[cookie.Name] = cookie
			}
			for name, sameSite := range test.cookies {
				cookie, ok := got[name]
				if !ok {
					t.Errorf("missing clearing Set-Cookie for %q; headers = %v", name, response.Header.Values("Set-Cookie"))
					continue
				}
				assertSessionCookie(t, cookie, sameSite, -1)
				if cookie.Value != "" {
					t.Errorf("cleared cookie %q value = %q, want empty", name, cookie.Value)
				}
				if !cookie.Expires.Before(time.Now()) {
					t.Errorf("cleared cookie %q expires = %v, want past", name, cookie.Expires)
				}
			}
		})
	}
}

func TestMutationOriginValidation(t *testing.T) {
	fixture := newContractFixture(t)
	tests := []struct {
		name   string
		origin string
		want   int
	}{
		{name: "missing", want: http.StatusForbidden},
		{name: "foreign", origin: "https://attacker.example", want: http.StatusForbidden},
		{name: "malformed", origin: "not-an-origin", want: http.StatusForbidden},
		{name: "allowed", origin: testOrigin, want: http.StatusNoContent},
		{name: "normalized", origin: "HTTPS://SCORE.EXAMPLE.TEST", want: http.StatusNoContent},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signout", nil)
			if test.origin != "" {
				request.Header.Set("Origin", test.origin)
			}
			recorder := httptest.NewRecorder()
			fixture.handler.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.want, recorder.Body.String())
			}
			if test.want == http.StatusForbidden {
				assertProblem(t, recorder, http.StatusForbidden, "origin_forbidden")
			}
		})
	}
}

func TestOAuthCookiesHaveProductionAttributesAndContractTTL(t *testing.T) {
	tests := []struct {
		name              string
		startPath         string
		callbackPath      string
		transientName     string
		sessionName       string
		sessionTTL        time.Duration
		seedContestant    bool
		wantCallbackClear []string
	}{
		{
			name: "contestant", startPath: "/api/v1/auth/discord", callbackPath: "/api/v1/auth/discord/callback",
			transientName: "oauth2-session", sessionName: "regalia-user-session", sessionTTL: service.ContestantTTL,
			seedContestant: true, wantCallbackClear: []string{"oauth2-session"},
		},
		{
			name: "signup", startPath: "/api/v1/auth/discord", callbackPath: "/api/v1/auth/discord/callback",
			transientName: "oauth2-session", sessionName: "signup-session", sessionTTL: service.SignupTTL,
			wantCallbackClear: []string{"oauth2-session"},
		},
		{
			name: "admin", startPath: "/api/v1/admin/auth/discord", callbackPath: "/api/v1/admin/auth/discord/callback",
			transientName: "admin-oauth2-session", sessionName: "admin-session", sessionTTL: service.AdminTTL,
			wantCallbackClear: []string{"admin-oauth2-session"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newContractFixture(t)
			if test.seedContestant {
				fixture.seedContestant(t)
			}
			start := fixture.request(t, http.MethodGet, test.startPath, "")
			if start.Code != http.StatusFound {
				t.Fatalf("start status = %d, want 302; body = %s", start.Code, start.Body.String())
			}
			startCookies := cookiesByName(start.Result())
			transient := startCookies[test.transientName]
			if transient == nil {
				t.Fatalf("missing %s cookie: %v", test.transientName, start.Result().Header.Values("Set-Cookie"))
			}
			assertSessionCookie(t, transient, http.SameSiteLaxMode, int(service.OAuthTTL/time.Second))

			location, err := url.Parse(start.Header().Get("Location"))
			if err != nil {
				t.Fatalf("parse authorization location: %v", err)
			}
			state := location.Query().Get("state")
			if state == "" {
				t.Fatal("authorization location omitted state")
			}
			callbackTarget := test.callbackPath + "?code=oauth-code&state=" + url.QueryEscape(state)
			callback := fixture.request(t, http.MethodGet, callbackTarget, "", transient)
			if callback.Code != http.StatusFound {
				t.Fatalf("callback status = %d, want 302; body = %s", callback.Code, callback.Body.String())
			}
			callbackCookies := cookiesByName(callback.Result())
			created := callbackCookies[test.sessionName]
			if created == nil {
				t.Fatalf("missing %s cookie: %v", test.sessionName, callback.Result().Header.Values("Set-Cookie"))
			}
			assertSessionCookie(t, created, http.SameSiteStrictMode, int(test.sessionTTL/time.Second))
			for _, name := range test.wantCallbackClear {
				cleared := callbackCookies[name]
				if cleared == nil || cleared.MaxAge >= 0 {
					t.Errorf("callback did not clear transient cookie %q: %v", name, callback.Result().Header.Values("Set-Cookie"))
				}
			}
		})
	}
}

func TestAdminAndContestantSessionGuardsAreSeparated(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	contestantCookie := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	adminCookie := fixture.createSession(t, session.Data{Kind: session.KindAdmin, AdminName: "operator"}, service.AdminTTL)

	tests := []struct {
		name    string
		path    string
		cookie  *http.Cookie
		want    int
		problem string
	}{
		{name: "admin anonymous", path: "/api/v1/admin/teams", want: http.StatusUnauthorized, problem: "invalid_session"},
		{name: "admin rejects contestant cookie", path: "/api/v1/admin/teams", cookie: contestantCookie, want: http.StatusUnauthorized, problem: "invalid_session"},
		{name: "admin accepts admin cookie", path: "/api/v1/admin/teams", cookie: adminCookie, want: http.StatusOK},
		{name: "contestant anonymous", path: "/api/v1/contestant/profile", want: http.StatusUnauthorized, problem: "invalid_session"},
		{name: "contestant rejects admin cookie", path: "/api/v1/contestant/profile", cookie: adminCookie, want: http.StatusUnauthorized, problem: "invalid_session"},
		{name: "contestant accepts contestant cookie", path: "/api/v1/contestant/profile", cookie: contestantCookie, want: http.StatusOK},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var cookies []*http.Cookie
			if test.cookie != nil {
				cookies = append(cookies, test.cookie)
			}
			recorder := fixture.request(t, http.MethodGet, test.path, "", cookies...)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.want, recorder.Body.String())
			}
			if test.problem != "" {
				assertProblem(t, recorder, test.want, test.problem)
			}
		})
	}
}

func TestSStateCallbackRequiresFixedBearerWithoutBrowserOrigin(t *testing.T) {
	fixture := newContractFixture(t)
	body := fmt.Sprintf(`{"event_id":"1091f318-87f0-4bb8-984d-688d0f413453","occurred_at":%q,"status":"DEPLOYING","message":null}`, testNow.Format(time.RFC3339))
	tests := []struct {
		name    string
		bearer  string
		want    int
		problem string
	}{
		{name: "missing", want: http.StatusUnauthorized, problem: "invalid_machine_token"},
		{name: "incorrect", bearer: "wrong-token", want: http.StatusUnauthorized, problem: "invalid_machine_token"},
		{name: "accepted", bearer: testSStateToken, want: http.StatusNotFound, problem: "resource_not_found"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/deployments/2/P1/1/events", strings.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			if test.bearer != "" {
				request.Header.Set("Authorization", "Bearer "+test.bearer)
			}
			recorder := httptest.NewRecorder()
			fixture.handler.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.want, recorder.Body.String())
			}
			assertProblem(t, recorder, test.want, test.problem)
		})
	}
}

func TestRateLimitUsesRFC9457AndRetryAfter(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	cookie := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)

	submit := func(body string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, "/api/v1/contestant/problems/P1/answers", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Origin", testOrigin)
		request.AddCookie(cookie)
		recorder := httptest.NewRecorder()
		fixture.handler.ServeHTTP(recorder, request)
		return recorder
	}
	first := submit(`{"body":"first answer"}`)
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want 201; body = %s", first.Code, first.Body.String())
	}
	second := submit(`{"body":"second answer"}`)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429; body = %s", second.Code, second.Body.String())
	}
	assertProblem(t, second, http.StatusTooManyRequests, "answer_rate_limited")
	if got, want := second.Header().Get("Retry-After"), fmt.Sprint(int(core.AnswerInterval/time.Second)); got != want {
		t.Errorf("Retry-After = %q, want %q", got, want)
	}
}

func TestDeploymentSSEStartsWithDatabaseSnapshotOnInitialAndReconnect(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	cookie := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	server := httptest.NewServer(fixture.handler)
	t.Cleanup(server.Close)

	for _, lastEventID := range []string{"", "previous-event-id"} {
		name := "initial"
		if lastEventID != "" {
			name = "reconnect"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/v1/contestant/problems/P1/deployments/stream", nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			request.AddCookie(cookie)
			if lastEventID != "" {
				request.Header.Set("Last-Event-ID", lastEventID)
			}
			response, err := server.Client().Do(request)
			if err != nil {
				t.Fatalf("open SSE: %v", err)
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want 200", response.StatusCode)
			}
			if got := response.Header.Get("Content-Type"); got != "text/event-stream" {
				t.Errorf("Content-Type = %q, want text/event-stream", got)
			}
			if got := response.Header.Get("Cache-Control"); got != "no-cache, no-store" {
				t.Errorf("Cache-Control = %q, want no-cache, no-store", got)
			}
			if got := response.Header.Get("X-Accel-Buffering"); got != "no" {
				t.Errorf("X-Accel-Buffering = %q, want no", got)
			}

			reader := bufio.NewReader(response.Body)
			event, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("read SSE event: %v", err)
			}
			id, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("read SSE id: %v", err)
			}
			data, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("read SSE data: %v", err)
			}
			if event != "event: snapshot\n" {
				t.Errorf("event line = %q, want snapshot", event)
			}
			if !strings.HasPrefix(id, "id: ") || strings.TrimSpace(strings.TrimPrefix(id, "id: ")) == "" {
				t.Errorf("id line = %q, want non-empty id", id)
			}
			var payload struct {
				Deployments []json.RawMessage `json:"deployments"`
			}
			if err := json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(data, "data: "))), &payload); err != nil {
				t.Fatalf("decode snapshot: %v; line = %q", err, data)
			}
			if payload.Deployments == nil || len(payload.Deployments) != 0 {
				t.Errorf("snapshot deployments = %v, want non-nil empty array", payload.Deployments)
			}
		})
	}
}

func TestGeneratedRouterRegistersEveryOpenAPIOperation(t *testing.T) {
	spec, err := api.GetSwagger()
	if err != nil {
		t.Fatalf("GetSwagger() error = %v", err)
	}
	want := make(map[string]struct{})
	for path, item := range spec.Paths.Map() {
		for method := range item.Operations() {
			want[strings.ToUpper(method)+" "+path] = struct{}{}
		}
	}

	router := chi.NewRouter()
	api.HandlerFromMux(api.Unimplemented{}, router)
	got := make(map[string]struct{})
	if err := chi.Walk(router, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		got[method+" "+route] = struct{}{}
		return nil
	}); err != nil {
		t.Fatalf("walk generated router: %v", err)
	}

	missing := setDifference(want, got)
	extra := setDifference(got, want)
	if len(want) != 61 {
		t.Errorf("OpenAPI operation count = %d, want 61", len(want))
	}
	if len(missing) != 0 || len(extra) != 0 {
		t.Errorf("generated route mismatch\nmissing: %v\nextra: %v", missing, extra)
	}
}

func assertSessionCookie(t *testing.T, cookie *http.Cookie, sameSite http.SameSite, maxAge int) {
	t.Helper()
	if !cookie.HttpOnly {
		t.Errorf("cookie %q is not HttpOnly", cookie.Name)
	}
	if !cookie.Secure {
		t.Errorf("cookie %q is not Secure", cookie.Name)
	}
	if cookie.Path != "/" {
		t.Errorf("cookie %q path = %q, want /", cookie.Name, cookie.Path)
	}
	if cookie.SameSite != sameSite {
		t.Errorf("cookie %q SameSite = %v, want %v", cookie.Name, cookie.SameSite, sameSite)
	}
	if maxAge < 0 {
		if cookie.MaxAge >= 0 {
			t.Errorf("cookie %q MaxAge = %d, want deletion", cookie.Name, cookie.MaxAge)
		}
		return
	}
	if cookie.MaxAge != maxAge {
		t.Errorf("cookie %q MaxAge = %d, want %d", cookie.Name, cookie.MaxAge, maxAge)
	}
	remaining := time.Until(cookie.Expires)
	if remaining < time.Duration(maxAge)*time.Second-time.Minute || remaining > time.Duration(maxAge)*time.Second+time.Minute {
		t.Errorf("cookie %q expires in %v, want approximately %v", cookie.Name, remaining, time.Duration(maxAge)*time.Second)
	}
}

func cookiesByName(response *http.Response) map[string]*http.Cookie {
	result := make(map[string]*http.Cookie)
	for _, cookie := range response.Cookies() {
		result[cookie.Name] = cookie
	}
	return result
}

func assertProblem(t *testing.T, recorder *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if got := recorder.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Errorf("Content-Type = %q, want application/problem+json", got)
	}
	var problem struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Code     string `json:"code"`
		Instance string `json:"instance"`
	}
	if err := json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v; body = %s", err, recorder.Body.String())
	}
	if problem.Status != status || problem.Code != code {
		t.Errorf("problem status/code = %d/%q, want %d/%q; body = %s", problem.Status, problem.Code, status, code, recorder.Body.String())
	}
	if problem.Type == "" || problem.Title == "" || problem.Detail == "" || problem.Instance == "" {
		t.Errorf("problem omitted RFC 9457 fields: %+v", problem)
	}
}

func setDifference(left, right map[string]struct{}) []string {
	result := make([]string, 0)
	for value := range left {
		if _, ok := right[value]; !ok {
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}
