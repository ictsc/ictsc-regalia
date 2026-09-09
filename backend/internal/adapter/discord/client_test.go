package discord

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestAuthorizationURLUsesPKCEAndExpectedScopes(t *testing.T) {
	client := newTestClient(t, "http://discord.invalid", nil)

	contestant, err := url.Parse(client.AuthorizationURL("state-value", "challenge-value", "https://score.example/callback", false))
	if err != nil {
		t.Fatal(err)
	}
	assertQueryValue(t, contestant.Query(), "scope", "identify")
	assertQueryValue(t, contestant.Query(), "state", "state-value")
	assertQueryValue(t, contestant.Query(), "code_challenge", "challenge-value")
	assertQueryValue(t, contestant.Query(), "code_challenge_method", "S256")

	admin, err := url.Parse(client.AuthorizationURL("admin-state", "admin-challenge", "https://score.example/admin/callback", true))
	if err != nil {
		t.Fatal(err)
	}
	assertQueryValue(t, admin.Query(), "scope", "identify guilds.members.read")
}

func TestExchangeFetchesUserAndAdminGuildMember(t *testing.T) {
	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.URL.Path {
		case "/oauth2/token":
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if r.Form.Get("grant_type") != "authorization_code" || r.Form.Get("code") != "the-code" ||
				r.Form.Get("code_verifier") != "the-verifier" || r.Form.Get("redirect_uri") != "https://score.example/admin/callback" ||
				r.Form.Get("client_id") != "client-id" || r.Form.Get("client_secret") != "client-secret" {
				t.Fatalf("unexpected token form: %#v", r.Form)
			}
			writeJSON(t, w, map[string]any{
				"access_token": "access-token", "token_type": "Bearer", "expires_in": 3600,
				"scope": "identify guilds.members.read", "documented_extra": true,
			})
		case "/api/users/@me":
			assertBearer(t, r, "access-token")
			writeJSON(t, w, map[string]any{
				"id": "123456789012345678", "username": "operator", "global_name": "ICTSC Operator", "avatar": nil,
			})
		case "/api/users/@me/guilds/987654321098765432/member":
			assertBearer(t, r, "access-token")
			writeJSON(t, w, map[string]any{
				"roles": []string{"333333333333333333", "222222222222222222", "333333333333333333"},
				"nick":  "staff",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, server.Client())
	result, err := client.Exchange(context.Background(), "the-code", "the-verifier", "https://score.example/admin/callback", true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Identity.ID != "123456789012345678" || result.Identity.DisplayName != "ICTSC Operator" {
		t.Fatalf("unexpected identity: %#v", result.Identity)
	}
	if result.GuildID != "987654321098765432" || !reflect.DeepEqual(result.RoleIDs, []string{"222222222222222222", "333333333333333333"}) {
		t.Fatalf("unexpected guild result: %#v", result)
	}
	if !client.HasAllowedAdminRole(result.RoleIDs) {
		t.Fatal("expected configured admin role to be accepted")
	}
	if err := client.RequireAllowedAdminRole(result.RoleIDs); err != nil {
		t.Fatalf("allowed roles were rejected: %v", err)
	}
	if err := client.RequireAllowedAdminRole([]string{"444444444444444444"}); !errors.Is(err, ErrAdminRoleRequired) {
		t.Fatalf("unallowed roles error = %v, want ErrAdminRoleRequired", err)
	}
	wantCalls := []string{"POST /oauth2/token", "GET /api/users/@me", "GET /api/users/@me/guilds/987654321098765432/member"}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("calls = %#v, want %#v", calls, wantCalls)
	}
}

func TestExchangeContestantDoesNotFetchGuildMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/token":
			writeJSON(t, w, map[string]any{"access_token": "access-token", "token_type": "Bearer"})
		case "/api/users/@me":
			writeJSON(t, w, map[string]any{"id": "123456789012345678", "username": "contestant", "global_name": nil})
		default:
			t.Fatalf("unexpected request: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, server.Client())
	result, err := client.Exchange(context.Background(), "the-code", "the-verifier", "https://score.example/callback", false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Identity.DisplayName != "contestant" || result.GuildID != "" || len(result.RoleIDs) != 0 {
		t.Fatalf("unexpected contestant result: %#v", result)
	}
}

func TestGuildMembershipNotFoundHasStableError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth2/token":
			writeJSON(t, w, map[string]any{"access_token": "access-token", "token_type": "Bearer"})
		case "/api/users/@me":
			writeJSON(t, w, map[string]any{"id": "123456789012345678", "username": "operator"})
		default:
			http.Error(w, `{"message":"Unknown Guild"}`, http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, server.Client())
	_, err := client.Exchange(context.Background(), "the-code", "the-verifier", "https://score.example/admin/callback", true)
	if !errors.Is(err, ErrGuildMembership) {
		t.Fatalf("error = %v, want ErrGuildMembership", err)
	}
}

func TestPKCEGeneration(t *testing.T) {
	verifier, challenge, err := GeneratePKCE()
	if err != nil {
		t.Fatal(err)
	}
	if len(verifier) != 43 || len(challenge) != 43 || strings.ContainsAny(verifier+challenge, "+/=") {
		t.Fatalf("unexpected PKCE pair: verifier=%q challenge=%q", verifier, challenge)
	}
	if challenge != S256Challenge(verifier) {
		t.Fatal("challenge is not S256(verifier)")
	}
}

func newTestClient(t *testing.T, baseURL string, httpClient *http.Client) *Client {
	t.Helper()
	client, err := New(Config{
		ClientID: "client-id", ClientSecret: "client-secret", AdminGuildID: "987654321098765432",
		AllowedAdminRoleIDs:   []string{"222222222222222222"},
		AuthorizationEndpoint: baseURL + "/oauth2/authorize", TokenEndpoint: baseURL + "/oauth2/token",
		APIBaseURL: baseURL + "/api", HTTPClient: httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func assertQueryValue(t *testing.T, values url.Values, key, want string) {
	t.Helper()
	if got := values.Get(key); got != want {
		t.Fatalf("query %s = %q, want %q", key, got, want)
	}
}

func assertBearer(t *testing.T, request *http.Request, token string) {
	t.Helper()
	if request.Header.Get("Authorization") != "Bearer "+token {
		t.Fatalf("Authorization = %q", request.Header.Get("Authorization"))
	}
}

func writeJSON(t *testing.T, writer http.ResponseWriter, value any) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		t.Fatal(err)
	}
}
