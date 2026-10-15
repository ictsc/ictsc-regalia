package github

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVerifyActionsToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 30, 1, 2, 3, 0, time.UTC)
	var discoveryCalls, jwksCalls int
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			discoveryCalls++
			writeTestJSON(t, w, map[string]any{"issuer": server.URL, "jwks_uri": server.URL + "/.well-known/jwks", "extra": true})
		case "/.well-known/jwks":
			jwksCalls++
			writeTestJSON(t, w, map[string]any{"keys": []any{jwkForKey("actions-key", &key.PublicKey)}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := newGitHubTestClient(t, server.URL, server.Client(), now)
	claims := validClaims(server.URL, now)
	claims["actor"] = "ictsc-bot"
	token := signJWT(t, key, "actions-key", claims)
	got, err := client.VerifyActionsToken(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}
	if got.RepositoryID != "987654321" || got.Repository != "ictsc/content" || got.Ref != "refs/heads/main" ||
		got.WorkflowRef != "ictsc/content/.github/workflows/publish.yml@refs/heads/main" || got.Subject != "repo:ictsc/content:ref:refs/heads/main" {
		t.Fatalf("unexpected claims: %#v", got)
	}
	if _, err := client.VerifyActionsToken(context.Background(), token); err != nil {
		t.Fatal(err)
	}
	if discoveryCalls != 1 || jwksCalls != 1 {
		t.Fatalf("OIDC documents were not cached: discovery=%d jwks=%d", discoveryCalls, jwksCalls)
	}
}

func TestVerifyActionsTokenRejectsBoundClaimMismatch(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 30, 1, 2, 3, 0, time.UTC)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			writeTestJSON(t, w, map[string]any{"issuer": server.URL, "jwks_uri": server.URL + "/jwks"})
		case "/jwks":
			writeTestJSON(t, w, map[string]any{"keys": []any{jwkForKey("actions-key", &key.PublicKey)}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := newGitHubTestClient(t, server.URL, server.Client(), now)

	tests := map[string]func(map[string]any){
		"issuer":        func(claims map[string]any) { claims["iss"] = "https://issuer.invalid" },
		"audience":      func(claims map[string]any) { claims["aud"] = "wrong-audience" },
		"repository id": func(claims map[string]any) { claims["repository_id"] = "123" },
		"repository":    func(claims map[string]any) { claims["repository"] = "attacker/content" },
		"ref":           func(claims map[string]any) { claims["ref"] = "refs/heads/other" },
		"workflow": func(claims map[string]any) {
			claims["job_workflow_ref"] = "ictsc/content/.github/workflows/other.yml@refs/heads/main"
		},
		"expired": func(claims map[string]any) { claims["exp"] = now.Add(-time.Minute).Unix() },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			claims := validClaims(server.URL, now)
			mutate(claims)
			_, err := client.VerifyActionsToken(context.Background(), signJWT(t, key, "actions-key", claims))
			if !errors.Is(err, ErrInvalidActionsToken) {
				t.Fatalf("error = %v, want ErrInvalidActionsToken", err)
			}
		})
	}
}

func TestVerifyActionsTokenRejectsDuplicateClaim(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"key"}`))
	claims := base64.RawURLEncoding.EncodeToString([]byte(`{"iss":"first","iss":"second"}`))
	token := header + "." + claims + ".c2lnbmF0dXJl"

	client := &Client{}
	_, err := client.VerifyActionsToken(context.Background(), token)
	if !errors.Is(err, ErrInvalidActionsToken) {
		t.Fatalf("error = %v, want ErrInvalidActionsToken", err)
	}
}

func newGitHubTestClient(t *testing.T, serverURL string, httpClient *http.Client, now time.Time) *Client {
	t.Helper()
	client, err := New(Config{
		APIBaseURL: serverURL, Issuer: serverURL, DiscoveryURL: serverURL + "/.well-known/openid-configuration",
		Audience: "ictsc-regalia-content", RepositoryID: "987654321", Repository: "ictsc/content",
		Ref: "refs/heads/main", WorkflowRef: "ictsc/content/.github/workflows/publish.yml@refs/heads/main",
		HTTPClient: httpClient, Now: func() time.Time { return now }, AllowInsecureHTTP: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func validClaims(issuer string, now time.Time) map[string]any {
	return map[string]any{
		"iss": issuer, "aud": []string{"another-audience", "ictsc-regalia-content"},
		"sub": "repo:ictsc/content:ref:refs/heads/main", "iat": now.Add(-time.Minute).Unix(),
		"nbf": now.Add(-time.Minute).Unix(), "exp": now.Add(5 * time.Minute).Unix(),
		"repository_id": "987654321", "repository": "ictsc/content", "ref": "refs/heads/main",
		"job_workflow_ref": "ictsc/content/.github/workflows/publish.yml@refs/heads/main",
	}
}

func signJWT(t *testing.T, key *rsa.PrivateKey, keyID string, claims map[string]any) string {
	t.Helper()
	headerBytes, err := json.Marshal(map[string]any{"alg": "RS256", "kid": keyID, "typ": "JWT", "x5t": "documented-extra"})
	if err != nil {
		t.Fatal(err)
	}
	claimBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	header := base64.RawURLEncoding.EncodeToString(headerBytes)
	payload := base64.RawURLEncoding.EncodeToString(claimBytes)
	digest := sha256.Sum256([]byte(header + "." + payload))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	return header + "." + payload + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func jwkForKey(keyID string, key *rsa.PublicKey) map[string]any {
	return map[string]any{
		"kty": "RSA", "use": "sig", "kid": keyID, "alg": "RS256",
		"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}

func writeTestJSON(t *testing.T, writer http.ResponseWriter, value any) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		t.Fatal(err)
	}
}
