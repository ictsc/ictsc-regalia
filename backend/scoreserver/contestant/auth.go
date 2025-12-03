package contestant

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
	"github.com/ictsc/ictsc-regalia/backend/pkg/ratelimiter"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/discord"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

type (
	AuthHandler struct {
		handler http.Handler
		once    sync.Once

		TrustProxy          bool
		ExternalURL         *url.URL
		DiscordOAuth2Config oauth2.Config
		RateLimiter         ratelimiter.RateLimiter

		DiscordCallbackEffect DiscordCallbackEffect
		SignUpEffect          SignUpEffect
	}
	DiscordCallbackEffect interface {
		domain.DiscordIdentityGetter
		domain.DiscordLinkedUserGetter
	}
)

type oauthErrorCode string

const (
	oauthErrorInvalidRequest         oauthErrorCode = "invalid_request"
	oauthErrorInvalidClient          oauthErrorCode = "invalid_client"
	oauthErrorInvalidGrant           oauthErrorCode = "invalid_grant"
	oauthErrorUnauthorizedClient     oauthErrorCode = "unauthorized_client"
	oauthErrorUnsupportedGrantType   oauthErrorCode = "unsupported_grant_type"
	oauthErrorInvalidScope           oauthErrorCode = "invalid_scope"
	oauthErrorServerError            oauthErrorCode = "server_error"
	oauthErrorTemporarilyUnavailable oauthErrorCode = "temporarily_unavailable"
)
const (
	authCookieAge   = 10 * time.Minute
	signUpCookieAge = 10 * time.Minute
	userCookieAge   = 3 * 24 * time.Hour
)

func newAuthHandler(cfg config.ContestantAuth, repo *pg.Repository, rateLimiter ratelimiter.RateLimiter) *AuthHandler {
	return &AuthHandler{
		TrustProxy:  cfg.TrustProxy,
		ExternalURL: cfg.ExternalURL,
		DiscordOAuth2Config: oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
			Scopes: []string{"identify"},
		},
		RateLimiter: rateLimiter,

		DiscordCallbackEffect: struct {
			*discord.UserClient
			*pg.Repository
		}{
			UserClient: &discord.UserClient{HTTPClient: otelhttp.DefaultClient},
			Repository: repo,
		},
		SignUpEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) SignUpTxEffect { return rt }),
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("GET /auth/signin",
			otelhttp.WithRouteTag("/auth/signin", http.HandlerFunc(h.handleSignIn)))
		mux.Handle("GET /auth/callback",
			otelhttp.WithRouteTag("/auth/callback", http.HandlerFunc(h.handleCallback)))
		mux.Handle("POST /auth/signup",
			otelhttp.WithRouteTag("/auth/signup", http.HandlerFunc(h.handleSignUp)))
		mux.Handle("POST /auth/signout",
			otelhttp.WithRouteTag("/auth/signout", http.HandlerFunc(h.handleSignOut)))
		h.handler = mux
	})
	h.handler.ServeHTTP(w, r)
}

func (h *AuthHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	nextPathURL := &url.URL{Path: "/"}
	if nextPath := r.URL.Query().Get("next"); nextPath != "" {
		url, err := parsePath(nextPath)
		if err == nil {
			nextPathURL = url
		}
	}

	sess := h.generateOAuth2Session(nextPathURL)
	h.discordAuthCodeURLRedirect(w, r, sess)
}

func (h *AuthHandler) generateOAuth2Session(nextPath *url.URL) *session.OAuth2Session {
	return &session.OAuth2Session{
		State:    oauth2.GenerateVerifier(),
		Verifier: oauth2.GenerateVerifier(),
		NextPath: nextPath,
	}
}

func (h *AuthHandler) externalURL(r *http.Request) *url.URL {
	url := &url.URL{}

	// From URL
	if h.ExternalURL != nil {
		if url.Scheme == "" && h.ExternalURL.Scheme != "" {
			url.Scheme = h.ExternalURL.Scheme
		}
		if url.Host == "" && h.ExternalURL.Host != "" {
			url.Host = h.ExternalURL.Host
		}
		if url.Path == "" && h.ExternalURL.Path != "" {
			url.Path = h.ExternalURL.Path
		}
	}
	// From Proxy Headers
	if h.TrustProxy {
		if proto := r.Header.Get("X-Forwarded-Proto"); url.Scheme == "" && (proto == "http" || proto == "https") {
			url.Scheme = proto
		}
		if host := r.Header.Get("X-Forwarded-Host"); url.Host == "" && host != "" {
			// Validate host doesn't contain path separators or other malicious content
			if !strings.ContainsAny(host, "/\\@") {
				url.Host = host
			}
		}
	}
	// From Request
	if url.Scheme == "" {
		if r.TLS != nil {
			url.Scheme = "https"
		} else {
			url.Scheme = "http"
		}
	}
	if url.Host == "" {
		url.Host = r.Host
	}

	return url
}

func (h *AuthHandler) discordAuthCodeURLRedirect(w http.ResponseWriter, r *http.Request, sess *session.OAuth2Session) {
	oauthCfg := h.DiscordOAuth2Config
	oauthCfg.RedirectURL = h.externalURL(r).JoinPath("./auth/callback").String()

	authCodeURL := oauthCfg.AuthCodeURL(
		sess.State,
		oauth2.S256ChallengeOption(sess.Verifier),
	)
	if err := session.OAuth2SessionStore.Write(r, w, sess, h.oauth2SessionOption(r)); err != nil {
		slog.ErrorContext(r.Context(), "failed to write oauth2 session", "error", errors.WithStack(err))
		h.errorRedirect(w, r, sess.NextPath, oauthErrorServerError)
		return
	}

	http.Redirect(w, r, authCodeURL, http.StatusFound)
}

func (h *AuthHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	// Get session
	oauthSess, err := session.OAuth2SessionStore.Get(r.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.errorRedirect(w, r, &url.URL{Path: "/"}, oauthErrorInvalidRequest)
			return
		}
		slog.ErrorContext(r.Context(), "failed to get session", "error", errors.WithStack(err))
		h.errorRedirect(w, r, &url.URL{Path: "/"}, oauthErrorServerError)
		return
	}

	// Delete session
	if err := session.OAuth2SessionStore.Write(r, w, nil, h.oauth2SessionOption(r)); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete session", "error", errors.WithStack(err))
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	// Exchange auth code
	if r.URL.Query().Get("state") != oauthSess.State {
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorInvalidRequest)
		return
	}

	if queryErr := r.URL.Query().Get("error"); queryErr != "" {
		desc := r.URL.Query().Get("error_description")
		slog.InfoContext(r.Context(), "failed to get discord callback", "error", queryErr, "error_description", desc)
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorCode(queryErr))
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorInvalidRequest)
		return
	}

	oauthCfg := h.DiscordOAuth2Config
	oauthCfg.RedirectURL = h.externalURL(r).JoinPath("./auth/callback").String()

	token, err := oauthCfg.Exchange(
		context.WithValue(r.Context(), oauth2.HTTPClient, otelhttp.DefaultClient),
		code, oauth2.VerifierOption(oauthSess.Verifier),
	)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to exchange code", "error", errors.WithStack(err))
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	identity, err := domain.GetDiscordIdentity(r.Context(), h.DiscordCallbackEffect, token.AccessToken)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get discord identity", "error", err)
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	user, err := identity.ID().User(r.Context(), h.DiscordCallbackEffect)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		slog.ErrorContext(r.Context(), "failed to get discord linked user", "error", err)
		h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	if user != nil {
		// 登録済み
		userSess := &session.UserSession{UserID: uuid.UUID(user.ID())}
		if err := session.UserSessionStore.Write(r, w, userSess, h.userSessionOption(r)); err != nil {
			slog.ErrorContext(r.Context(), "failed to write user session", "error", err)
			h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorServerError)
			return
		}
	} else {
		signUpSess := &session.SignUpSession{Discord: identity.Data()}
		if err := session.SignUpSessionStore.Write(r, w, signUpSess, h.signUpSessionOption(r)); err != nil {
			slog.ErrorContext(r.Context(), "failed to write signup session", "error", err)
			h.errorRedirect(w, r, oauthSess.NextPath, oauthErrorServerError)
			return
		}
	}

	http.Redirect(w, r, oauthSess.NextPath.String(), http.StatusFound)
}

func (h *AuthHandler) handleSignOut(w http.ResponseWriter, r *http.Request) {
	if err := session.UserSessionStore.Write(r, w, nil, h.userSessionOption(r)); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete user session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := session.SignUpSessionStore.Write(r, w, nil, h.signUpSessionOption(r)); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete signup session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) oauth2SessionOption(r *http.Request) *sessions.Options {
	opt := h.sessionOption(r)
	opt.MaxAge = int(authCookieAge.Seconds())
	opt.Path = h.externalURL(r).JoinPath("./auth").Path
	opt.SameSite = http.SameSiteLaxMode
	return opt
}

func (h *AuthHandler) userSessionOption(r *http.Request) *sessions.Options {
	opt := h.sessionOption(r)
	opt.MaxAge = int(userCookieAge.Seconds())
	return opt
}

func (h *AuthHandler) signUpSessionOption(r *http.Request) *sessions.Options {
	opt := h.sessionOption(r)
	opt.MaxAge = int(signUpCookieAge.Seconds())
	return opt
}

func (h *AuthHandler) sessionOption(r *http.Request) *sessions.Options {
	externalURL := h.externalURL(r)
	path := externalURL.Path
	if path == "" {
		path = "/"
	}
	return &sessions.Options{
		Path:     path,
		Secure:   externalURL.Scheme == "https",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (h *AuthHandler) errorRedirect(w http.ResponseWriter, r *http.Request, nextPath *url.URL, errCode oauthErrorCode) {
	frag := url.Values{}
	frag.Set("oauth-error", string(errCode))
	nextPath.Fragment = frag.Encode()
	http.Redirect(w, r, nextPath.String(), http.StatusFound)
}

func parsePath(path string) (*url.URL, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("path must start with /")
	}

	parsedURL, err := url.Parse(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse path")
	}

	if parsedURL.RawQuery != "" {
		return nil, errors.New("next path must not contain query")
	}

	if parsedURL.Host != "" {
		return nil, errors.New("next path must not contain host")
	}

	return parsedURL, nil
}
