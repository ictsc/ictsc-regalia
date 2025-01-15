package auth

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/sessions"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

type AuthOption interface {
	HandlerOption(h *handler)
}

type nowOption struct {
	now func() time.Time
}

func (n nowOption) HandlerOption(h *handler) {
	h.now = n.now
}

func WithNow(now func() time.Time) AuthOption {
	return nowOption{now: now}
}

func NewHandler(cfg config.ContestantAuth, sessStore sessions.Store, opts ...AuthOption) http.Handler {
	handler := &handler{
		sessStore: sessStore,
		now:       time.Now,

		baseURL: cfg.BaseURL,
	}
	for _, opt := range opts {
		opt.HandlerOption(handler)
	}

	handler.discordOAuth2Cfg = oauth2.Config{
		ClientID:     cfg.DiscordClientID,
		ClientSecret: cfg.DiscordClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		RedirectURL: handler.baseURL.JoinPath("./auth/discord/callback").String(),
		Scopes:      []string{"identify"},
	}

	mux := http.NewServeMux()

	mux.Handle("GET /auth/discord", otelhttp.WithRouteTag("/auth/discord", http.HandlerFunc(handler.handleDiscordLogin)))
	mux.Handle("GET /auth/discord/callback", otelhttp.WithRouteTag("/auth/discord/callback", http.HandlerFunc(handler.handleDiscordCallback)))

	return mux
}

type handler struct {
	sessStore sessions.Store
	now       func() time.Time

	baseURL          url.URL
	discordOAuth2Cfg oauth2.Config
}

func (h *handler) handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	nextPath := "/"

	authState := generateAuthState()
	codeURL := h.discordOAuth2Cfg.AuthCodeURL(authState.State, oauth2.S256ChallengeOption(authState.Verifier))

	sess, err := h.sessStore.Get(r, authSessionKey)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get session", "error", errors.WithStack(err))
		sendError(w, r, nextPath)
		return
	}
	sess.Options.MaxAge = 10 * 60
	authState.save(sess)
	if err := h.sessStore.Save(r, w, sess); err != nil {
		slog.ErrorContext(r.Context(), "failed to save session", "error", errors.WithStack(err))
		sendError(w, r, nextPath)
		return
	}

	http.Redirect(w, r, codeURL, http.StatusFound)
}

func (h *handler) handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	nextPath := "/"

	authSess, err := h.sessStore.Get(r, authSessionKey)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get session", "error", errors.WithStack(err))
		sendError(w, r, nextPath)
		return
	}
	if authSess.IsNew {
		// no sesson
		sendError(w, r, nextPath)
		return
	}
	// clear session
	authSess.Options.MaxAge = -1
	if err := h.sessStore.Save(r, w, authSess); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete session", "error", errors.WithStack(err))
		sendError(w, r, nextPath)
		return
	}

	var authState authState
	authState.load(authSess)

	queryState := r.URL.Query().Get("state")
	if queryState != authState.State {
		// state unmatch
		sendError(w, r, nextPath)
		return
	}

	queryError := r.URL.Query().Get("error")
	if queryError != "" {
		description := r.URL.Query().Get("error_description")
		slog.InfoContext(r.Context(), "discord auth error", "error", queryError, "description", description)
		sendError(w, r, nextPath) //TODO: より詳細なエラーを報告する
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		sendError(w, r, nextPath)
		return
	}

	token, err := h.discordOAuth2Cfg.Exchange(r.Context(), code, oauth2.VerifierOption(authState.Verifier))
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to exchange token", "error", errors.WithStack(err))
		sendError(w, r, nextPath)
		return
	}

	slog.InfoContext(r.Context(), "discord auth success", "token", token)
}

func sendError(w http.ResponseWriter, r *http.Request, path string) {
	path += "#auth_error"
	http.Redirect(w, r, path, http.StatusFound)
}

const authSessionKey = "auth-session"

type authState struct {
	State    string
	Verifier string
}

func generateAuthState() *authState {
	return &authState{
		State:    oauth2.GenerateVerifier(),
		Verifier: oauth2.GenerateVerifier(),
	}
}

func (s *authState) load(sess *sessions.Session) {
	if stateAny, ok := sess.Values["state"]; ok {
		if state, ok := stateAny.(string); ok {
			s.State = state
		}
	}
	if verifierAny, ok := sess.Values["verifier"]; ok {
		if verifier, ok := verifierAny.(string); ok {
			s.Verifier = verifier
		}
	}
}

func (s *authState) save(sess *sessions.Session) {
	if s.State != "" {
		sess.Values["state"] = s.State
	}
	if s.Verifier != "" {
		sess.Values["verifier"] = s.Verifier
	}
}
