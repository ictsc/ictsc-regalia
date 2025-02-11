package session

import (
	"context"
	"encoding/gob"
	"net/http"
	"net/url"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

type (
	OAuth2Session struct {
		State    string
		Verifier string
		NextPath *url.URL
	}
	SignUpSession struct {
		Discord *domain.DiscordIdentityData
	}
	UserSession struct {
		UserID uuid.UUID
	}

	sessionCtxKey struct{}
	sessionCtx    struct {
		store    sessions.Store
		registry *sessions.Registry
	}
)

//nolint:gochecknoinits // gob.Register は init で呼び出す必要がある
func init() {
	gob.Register(&OAuth2Session{})
	gob.Register(&SignUpSession{})
	gob.Register(&domain.DiscordIdentityData{})
	gob.Register(&UserSession{})
}

func NewHandler(store sessions.Store) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			registry := sessions.GetRegistry(r)
			r = r.WithContext(context.WithValue(r.Context(), sessionCtxKey{}, &sessionCtx{
				store:    store,
				registry: registry,
			}))
			handler.ServeHTTP(w, r)
		})
	}
}

var (
	OAuth2SessionStore = &SessionStore[*OAuth2Session]{key: "oauth2-session"}
	SignUpSessionStore = &SessionStore[*SignUpSession]{key: "signup-session"}
	UserSessionStore   = &SessionStore[*UserSession]{key: "user-session"}
)

type SessionStore[V comparable] struct {
	key string
}

func (s *SessionStore[V]) Get(ctx context.Context) (V, error) {
	var zero V

	sessCtx, err := getSessCtx(ctx)
	if err != nil {
		return zero, err
	}

	sess, err := sessCtx.getSession(s.key)
	if err != nil {
		return zero, err
	}

	valAny, ok := sess.Values[0]
	if !ok {
		return zero, domain.NewNotFoundError("session value", nil)
	}
	val, ok := valAny.(V)
	if !ok {
		return zero, errors.New("session value has invalid type")
	}
	return val, nil
}

func (s *SessionStore[V]) Write(r *http.Request, w http.ResponseWriter, val V, options *sessions.Options) error {
	ctx := r.Context()

	sessCtx, err := getSessCtx(ctx)
	if err != nil {
		return err
	}

	sess, err := sessCtx.getSession(s.key)
	if err != nil {
		return err
	}

	var zero V
	if val != zero {
		sess.Values[0] = val
		if options != nil {
			sess.Options = options
		}
	} else {
		sess.Options.MaxAge = -1
	}

	if err := sess.Save(r, w); err != nil {
		return errors.Wrap(err, "failed to save session")
	}
	return nil
}

func getSessCtx(ctx context.Context) (*sessionCtx, error) {
	val, ok := ctx.Value(sessionCtxKey{}).(*sessionCtx)
	if !ok {
		return nil, errors.New("session store is not set")
	}
	return val, nil
}

func (c *sessionCtx) getSession(key string) (*sessions.Session, error) {
	sess, err := c.registry.Get(c.store, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}
	return sess, nil
}
