package contestant

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/httputil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/pkg/ratelimiter"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/connectdomain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(ctx context.Context, cfg config.ContestantAPI, db *sqlx.DB, rdb redis.UniversalClient) (http.Handler, error) {
	repo := pg.NewRepository(db)
	sessionStore, err := redisstore.NewRedisStore(ctx, rdb)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session store")
	}
	sessionStore.KeyPrefix("contestant-session:")
	rateLimiter := ratelimiter.NewRedisRateLimiter(rdb, "contestant-rate-limiter:")

	interceptors := []connect.Interceptor{
		connectutil.NewOtelInterceptor(),
		connectdomain.NewErrorInterceptor(),
		connectdomain.NewLoggingInterceptor(),
	}

	mux := http.NewServeMux()

	mux.Handle("/auth/", otelhttp.NewHandler(newAuthHandler(cfg.Auth, repo, rateLimiter), "auth"))

	mux.Handle(contestantv1connect.NewViewerServiceHandler(
		newViewerServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewProfileServiceHandler(
		newProfileServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewContestServiceHandler(
		newContestServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewProblemServiceHandler(
		newProblemServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewAnswerServiceHandler(
		newAnswerServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewNoticeServiceHandler(
		newNoticeServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))

	handler := http.Handler(mux)
	handler = session.NewHandler(sessionStore)(handler)
	handler = httputil.CSRFMiddleware(cfg.AllowedOrigins)(handler)
	handler = h2c.NewHandler(handler, &http2.Server{})

	return handler, nil
}
