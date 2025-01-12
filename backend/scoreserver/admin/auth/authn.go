package auth

import (
	"context"
	"crypto/x509"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
)

type (
	Viewer struct {
		Name   string
		Groups []string
	}
	viewerKey struct{}

	HTTPAuthenticator interface {
		HandleRequest(req *http.Request) (*Viewer, error)
	}
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

func GetViewer(ctx context.Context) Viewer {
	viewer, ok := ctx.Value(viewerKey{}).(*Viewer)
	if !ok || viewer == nil {
		return Viewer{
			Name:   "anonymous",
			Groups: []string{"system:unauthenticated"},
		}
	}
	return *viewer
}

func WithAuthn(handler http.Handler, authenticator HTTPAuthenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		viewer, err := authenticator.HandleRequest(r)
		if err != nil {
			if !errors.Is(err, ErrUnauthenticated) {
				slog.DebugContext(r.Context(), "failed to authenticate", "error", err)
				http.Error(w, "failed to authenticate", http.StatusInternalServerError)
				return
			}
		} else {
			viewer.Groups = append(viewer.Groups, "system:authenticated")
			ctx := context.WithValue(r.Context(), viewerKey{}, viewer)
			r = r.WithContext(ctx)
		}
		handler.ServeHTTP(w, r)
	})
}

type (
	JWTAuthenticator struct {
		issuers []*issuer
	}
	issuer struct {
		name      string
		verifier  *oidc.IDTokenVerifier
		nameKey   string
		groupKeys []string
	}
)

func NewJWTAuthenticator(ctx context.Context, cfg config.AdminAuthn) (*JWTAuthenticator, error) {
	issuers := make([]*issuer, 0, len(cfg.Issuers))
	for _, issuerCfg := range cfg.Issuers {
		iss, err := newIssuer(ctx, issuerCfg)
		if err != nil {
			return nil, err
		}
		issuers = append(issuers, iss)
	}

	return &JWTAuthenticator{issuers: issuers}, nil
}

func newIssuer(ctx context.Context, cfg config.Issuer) (*issuer, error) {
	if cfg.InsecureIssuerURL != "" {
		ctx = oidc.InsecureIssuerURLContext(ctx, cfg.InsecureIssuerURL)
	}

	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		transport = &http.Transport{}
	}

	if cfg.CAFile != "" {
		caPEM, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read CA file %s", cfg.CAFile)
		}

		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load system cert pool")
		}
		certPool = certPool.Clone()

		if !certPool.AppendCertsFromPEM(caPEM) {
			return nil, errors.New("failed to append CA to cert pool")
		}

		transport = transport.Clone()
		transport.TLSClientConfig.RootCAs = certPool
	}

	httpClient := &http.Client{Transport: transport}
	ctx = oidc.ClientContext(ctx, httpClient)

	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create provider for %s", cfg.Issuer)
	}
	verifier := provider.VerifierContext(ctx, &oidc.Config{ClientID: cfg.ClientID})

	return &issuer{
		name:      cfg.Issuer,
		verifier:  verifier,
		nameKey:   cfg.UsernameKey,
		groupKeys: cfg.GroupKeys,
	}, nil
}

func (a *JWTAuthenticator) HandleRequest(req *http.Request) (*Viewer, error) {
	authHeader := req.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrUnauthenticated
	}
	rawIDToken := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := req.Context()
	for _, iss := range a.issuers {
		idToken, err := iss.verifier.Verify(ctx, rawIDToken)
		if err != nil {
			err = errors.Wrapf(err, "for %s", iss.name)
			slog.DebugContext(ctx, "failed to verify id token", "error", err)
			continue
		}

		var claims map[string]any
		if err := idToken.Claims(&claims); err != nil {
			err = errors.Wrapf(err, "for %s", iss.name)
			slog.DebugContext(ctx, "failed to parse id token claims", "error", err)
			continue
		}

		viewerName := mapViewerName(idToken, claims, iss.nameKey)
		viewerGroups := mapViewerGroups(claims, iss.groupKeys)
		return &Viewer{
			Name:   viewerName,
			Groups: viewerGroups,
		}, nil
	}

	return nil, ErrUnauthenticated
}

func mapViewerName(idToken *oidc.IDToken, claims map[string]any, nameClaim string) string {
	defaultName := idToken.Subject

	if nameClaim == "" {
		nameClaim = "name"
	}
	nameRaw, ok := claims[nameClaim]
	if !ok {
		return defaultName
	}
	name, ok := nameRaw.(string)
	if !ok {
		return defaultName
	}
	return name
}

func mapViewerGroups(claims map[string]any, groupClaims []string) []string {
	groups := []string{}

	for _, claim := range groupClaims {
		groupRaw, ok := claims[claim]
		if !ok {
			continue
		}
		switch groupData := groupRaw.(type) {
		case string:
			groups = append(groups, groupData)
		case []any:
			for _, groupDataRaw := range groupData {
				if groupStr, ok := groupDataRaw.(string); ok {
					groups = append(groups, groupStr)
				}
			}
		}
	}

	return groups
}
