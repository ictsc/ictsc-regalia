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

		BaseURL               *url.URL
		DiscordOAuth2Config   oauth2.Config
		DiscordCallbackEffect DiscordCallbackEffect
	}
	DiscordCallbackEffect interface {
		domain.DiscordIdentityGetter
		domain.DiscordLinkedUserGetter
	}
	oauthErrorCode string
)

const (
	oauthErrorInvalidRequest         oauthErrorCode = "invalid_request"
	oauthErrorInvalidClient          oauthErrorCode = "invalid_client"
	oauthErrorInvalidGrant           oauthErrorCode = "invalid_grant"
	oauthErrorUnauthorizedClient     oauthErrorCode = "unauthorized_client"
	oauthErrorUnsupportedGrantType   oauthErrorCode = "unsupported_grant_type"
	oauthErrorInvalidScope           oauthErrorCode = "invalid_scope"
	oauthErrorServerError            oauthErrorCode = "server_error"
	oauthErrorTemporarilyUnavailable oauthErrorCode = "temporarily_unavailable"

	authCookieAge   = 10 * time.Minute
	signUpCookieAge = 10 * time.Minute
	userCookieAge   = 3 * 24 * time.Hour
)

func newAuthHandler(cfg config.ContestantAuth, repo *pg.Repository) *AuthHandler {
	return &AuthHandler{
		BaseURL: cfg.BaseURL,
		DiscordOAuth2Config: oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
			RedirectURL: cfg.BaseURL.JoinPath("./auth/callback").String(),
			Scopes:      []string{"identify"},
		},
		DiscordCallbackEffect: struct {
			*discord.UserClient
			*pg.Repository
		}{
			UserClient: &discord.UserClient{HTTPClient: otelhttp.DefaultClient},
			Repository: repo,
		},
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("GET /auth/signin",
			otelhttp.WithRouteTag("/auth/signin", http.HandlerFunc(h.handleSignIn)))
		mux.Handle("GET /auth/callback",
			otelhttp.WithRouteTag("/auth/callback", http.HandlerFunc(h.handleCallback)))
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

func (h *AuthHandler) discordAuthCodeURLRedirect(w http.ResponseWriter, r *http.Request, sess *session.OAuth2Session) {
	authCodeURL := h.DiscordOAuth2Config.AuthCodeURL(
		sess.State,
		oauth2.S256ChallengeOption(sess.Verifier),
	)
	if err := session.OAuth2SessionStore.Write(r, w, sess, &sessions.Options{
		MaxAge:   int(authCookieAge.Seconds()),
		Path:     h.BaseURL.JoinPath("./auth").Path,
		Secure:   h.BaseURL.Scheme == "https",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}); err != nil {
		slog.ErrorContext(r.Context(), "failed to write oauth2 session", "error", errors.WithStack(err))
		h.error(w, r, sess.NextPath, oauthErrorServerError)
		return
	}

	http.Redirect(w, r, authCodeURL, http.StatusFound)
}

func (h *AuthHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	// Get session
	oauthSess, err := session.OAuth2SessionStore.Get(r.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.error(w, r, &url.URL{Path: "/"}, oauthErrorInvalidRequest)
			return
		}
		slog.ErrorContext(r.Context(), "failed to get session", "error", errors.WithStack(err))
		h.error(w, r, &url.URL{Path: "/"}, oauthErrorServerError)
		return
	}

	// Delete session
	if err := session.OAuth2SessionStore.Write(r, w, nil, nil); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete session", "error", errors.WithStack(err))
		h.error(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	// Exchange auth code
	if r.URL.Query().Get("state") != oauthSess.State {
		h.error(w, r, oauthSess.NextPath, oauthErrorInvalidRequest)
		return
	}

	if queryErr := r.URL.Query().Get("error"); queryErr != "" {
		desc := r.URL.Query().Get("error_description")
		slog.InfoContext(r.Context(), "failed to get discord callback", "error", queryErr, "error_description", desc)
		h.error(w, r, oauthSess.NextPath, oauthErrorCode(queryErr))
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		h.error(w, r, oauthSess.NextPath, oauthErrorInvalidRequest)
		return
	}

	token, err := h.DiscordOAuth2Config.Exchange(
		context.WithValue(r.Context(), oauth2.HTTPClient, otelhttp.DefaultClient),
		code, oauth2.VerifierOption(oauthSess.Verifier),
	)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to exchange code", "error", errors.WithStack(err))
		h.error(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	identity, err := domain.GetDiscordIdentity(r.Context(), h.DiscordCallbackEffect, token.AccessToken)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get discord identity", "error", err)
		h.error(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	user, err := identity.ID().User(r.Context(), h.DiscordCallbackEffect)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		slog.ErrorContext(r.Context(), "failed to get discord linked user", "error", err)
		h.error(w, r, oauthSess.NextPath, oauthErrorServerError)
		return
	}

	path := h.BaseURL.Path
	if path == "" {
		path = "/"
	}
	sessOpt := &sessions.Options{
		Path:     path,
		Secure:   h.BaseURL.Scheme == "https",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	if user != nil {
		// 登録済み
		userSess := &session.UserSession{UserID: uuid.UUID(user.ID())}
		sessOpt.MaxAge = int(userCookieAge.Seconds())
		if err := session.UserSessionStore.Write(r, w, userSess, sessOpt); err != nil {
			slog.ErrorContext(r.Context(), "failed to write user session", "error", err)
			h.error(w, r, oauthSess.NextPath, oauthErrorServerError)
			return
		}
	} else {
		signUpSess := &session.SignUpSession{Discord: identity.Data()}
		sessOpt.MaxAge = int(signUpCookieAge.Seconds())
		if err := session.SignUpSessionStore.Write(r, w, signUpSess, sessOpt); err != nil {
			slog.ErrorContext(r.Context(), "failed to write signup session", "error", err)
			h.error(w, r, oauthSess.NextPath, oauthErrorServerError)
			return
		}
	}

	http.Redirect(w, r, oauthSess.NextPath.String(), http.StatusFound)
}

func (h *AuthHandler) error(w http.ResponseWriter, r *http.Request, nextPath *url.URL, errCode oauthErrorCode) {
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
