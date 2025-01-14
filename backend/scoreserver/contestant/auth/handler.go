package auth

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/discord", handler.discordLogin)

	return mux
}

type handler struct {
	sessStore sessions.Store
	now       func() time.Time

	baseURL   url.URL
	oauth2cfg oauth2.Config
}

func (h *handler) discordLogin(w http.ResponseWriter, r *http.Request) {
	redirectURI := h.baseURL.JoinPath("/auth/discord/callback")

	slog.InfoContext(r.Context(), "redirectURI", "redirectURI", redirectURI.String())
}
