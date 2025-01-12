package auth_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/oauth2"
)

func Test_AuthHandler(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		auth HTTPAuthenticatorFunc

		want           auth.Viewer
		wantStatusCode int
	}{
		"authenticated": {
			auth: func(*http.Request) (*auth.Viewer, error) {
				return &auth.Viewer{Name: "admin", Groups: []string{"admin"}}, nil
			},
			want: auth.Viewer{Name: "admin", Groups: []string{"admin", "system:authenticated"}},
		},
		"unauthenticated": {
			auth: func(*http.Request) (*auth.Viewer, error) {
				return nil, auth.ErrUnauthenticated
			},
			want: auth.Viewer{Name: "anonymous", Groups: []string{"system:unauthenticated"}},
		},
		"error": {
			auth: func(*http.Request) (*auth.Viewer, error) {
				return nil, errors.New("error")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := auth.WithAuthn(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					viewer := auth.GetViewer(r.Context())

					w.WriteHeader(http.StatusOK)
					if err := gob.NewEncoder(w).Encode(&viewer); err != nil {
						t.Errorf("Failed to encode response: %v", err)
					}
				}),
				tt.auth,
			)

			ctx := context.Background()
			req := httptest.NewRequestWithContext(ctx, "GET", "/", nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if t.Failed() {
				return
			}

			if tt.wantStatusCode != 0 && resp.Code != tt.wantStatusCode {
				t.Errorf("unexpected status code: %d", resp.Code)
			}
			if resp.Code != http.StatusOK {
				return
			}

			var viewer auth.Viewer
			if err := gob.NewDecoder(resp.Body).Decode(&viewer); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}
			if diff := cmp.Diff(tt.want, viewer); diff != "" {
				t.Errorf("unexpected viewer: %s", diff)
			}
		})
	}
}

type HTTPAuthenticatorFunc func(*http.Request) (*auth.Viewer, error)

func (f HTTPAuthenticatorFunc) HandleRequest(req *http.Request) (*auth.Viewer, error) {
	return f(req)
}

func Test_JWTAuth(t *testing.T) {
	t.Parallel()

	issuer := setupDex(t, dexConfig{})

	ctx := context.Background()
	token, err := retrieveIDTokenWithPasswordCredential(
		ctx, issuer,
		"app", "clientsecret",
		"admin@example.com", "password",
	)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	cases := map[string]struct {
		token  string
		config config.AdminAuthn

		want                *auth.Viewer
		wantUnauthenticated bool
	}{
		"valid": {
			token: token,
			config: config.AdminAuthn{
				Issuers: []config.Issuer{
					{
						Issuer:      issuer,
						ClientID:    "app",
						UsernameKey: "name",
						GroupKeys:   []string{"name"},
					},
				},
			},

			want: &auth.Viewer{
				Name:   "admin",
				Groups: []string{"admin"},
			},
		},
		"no issuer": {
			token:  token,
			config: config.AdminAuthn{},

			wantUnauthenticated: true,
		},
		"no token": {
			token:  token,
			config: config.AdminAuthn{},

			wantUnauthenticated: true,
		},
		"invalid audience": {
			token: token,
			config: config.AdminAuthn{
				Issuers: []config.Issuer{{Issuer: issuer, ClientID: "app2"}},
			},

			wantUnauthenticated: true,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tt.token != "" {
				req.Header.Add("Authorization", "Bearer "+tt.token)
			}

			authenticator, err := auth.NewJWTAuthenticator(ctx, tt.config)
			if err != nil {
				t.Fatalf("Failed to create authenticator: %v", err)
			}

			viewer, err := authenticator.HandleRequest(req)
			if errors.Is(err, auth.ErrUnauthenticated) != tt.wantUnauthenticated {
				t.Errorf("unexpected error: %v", err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.want, viewer); diff != "" {
				t.Errorf("unexpected viewer: %s", diff)
			}
		})
	}
}

func retrieveIDTokenWithPasswordCredential(
	ctx context.Context,
	issuer, clientID, clientSecret string,
	username, password string,
) (string, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return "", errors.Wrap(err, "failed to create provider")
	}

	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "profile", "email", "groups"},
		Endpoint:     provider.Endpoint(),
	}
	token, err := cfg.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("failed to get id_token")
	}

	if _, err := provider.VerifierContext(
		ctx, &oidc.Config{ClientID: clientID},
	).Verify(ctx, idToken); err != nil {
		return "", errors.Wrap(err, "failed to verify id_token")
	}

	return idToken, nil
}

type dexConfig struct{}

const (
	dexImage       = "dexidp/dex:v2.41.1"
	dexWebHTTPPort = "5556/tcp"
	dexConfigYAML  = `
issuer: %q

storage:
  type: sqlite3
  config:
    file: /var/dex/dex.db

web:
  http: 0.0.0.0:5556
telemetry:
  http: 0.0.0.0:5558

oauth2:
  passwordConnector: local

staticClients:
  - id: app
    name: "Example App"
    secret: clientsecret

enablePasswordDB: true
staticPasswords:
  - email: "admin@example.com"
    hash: "$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W"
    username: "admin"
    userID: "08a8684b-db88-4b73-90a9-3cd1661f5466"
`
)

func setupDex(tb testing.TB, _ dexConfig) string {
	tb.Helper()

	ctx := context.Background()

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: dexImage,
			Cmd: []string{
				"/bin/sh", "-c",
				"while [ ! -f /etc/dex/dex.yaml ]; do sleep 1; done; /usr/local/bin/dex serve /etc/dex/dex.yaml",
			},
			ExposedPorts: []string{dexWebHTTPPort, "5558/tcp"},
			Files: []testcontainers.ContainerFile{
				{
					ContainerFilePath: "/etc/dex/dex.yaml",
					FileMode:          int64(fs.ModePerm),
					Reader:            bytes.NewReader([]byte(dexConfigYAML)),
				},
			},
		},
		Logger:  testcontainers.TestLogger(tb),
		Started: true,
	})
	if err != nil {
		tb.Fatalf("Failed to create container: %v", err)
	}
	testcontainers.CleanupContainer(tb, ctr)

	host, err := ctr.Host(ctx)
	if err != nil {
		tb.Fatalf("Failed to get container host: %v", err)
	}

	httpPort, err := ctr.MappedPort(ctx, dexWebHTTPPort)
	if err != nil {
		tb.Fatalf("Failed to get container port: %v", err)
	}

	issuer := "http://" + net.JoinHostPort(host, httpPort.Port())
	if err := ctr.CopyToContainer(
		ctx, []byte(fmt.Sprintf(dexConfigYAML, issuer)), "/etc/dex/dex.yaml", int64(fs.ModePerm),
	); err != nil {
		tb.Fatalf("Failed to copy config: %v", err)
	}

	waitStrategy := wait.ForHTTP("/healthz/ready").WithPort(nat.Port("5558/tcp"))
	if err := waitStrategy.WaitUntilReady(ctx, ctr); err != nil {
		tb.Fatalf("Failed to wait for dex: %v", err)
	}

	return issuer
}
