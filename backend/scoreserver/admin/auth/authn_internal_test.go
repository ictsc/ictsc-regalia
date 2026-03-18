package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/coreos/go-oidc/v3/oidc"
)

type fakeVerifier struct {
	called int
}

func (v *fakeVerifier) Verify(context.Context, string) (*oidc.IDToken, error) {
	v.called++
	return nil, errors.New("verify failed")
}

func TestJWTAuthenticator_HandleRequest_UsesOnlyMatchingIssuer(t *testing.T) {
	t.Parallel()

	matchedVerifier := &fakeVerifier{}
	unmatchedVerifier := &fakeVerifier{}
	authenticator := JWTAuthenticator{
		issuers: map[string][]*issuer{
			"https://issuer.matched.example.com": {
				{name: "matched", verifier: matchedVerifier},
			},
			"https://issuer.unmatched.example.com": {
				{name: "unmatched", verifier: unmatchedVerifier},
			},
		},
	}

	token := jwtWithIssuer("https://issuer.matched.example.com")
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	_, err = authenticator.HandleRequest(req)
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expected unauthenticated error, got: %v", err)
	}

	if matchedVerifier.called != 1 {
		t.Fatalf("matched issuer verifier called %d times, want 1", matchedVerifier.called)
	}
	if unmatchedVerifier.called != 0 {
		t.Fatalf("unmatched issuer verifier called %d times, want 0", unmatchedVerifier.called)
	}
}

func jwtWithIssuer(issuer string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"iss":%q}`, issuer)))
	return header + "." + payload + ".signature"
}
