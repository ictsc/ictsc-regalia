package contestant

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/pkg/ratelimiter"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/connectdomain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	sloghttp "github.com/samber/slog-http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(
	ctx context.Context,
	cfg config.ContestantAPI,
	db *sqlx.DB,
	rdb redis.UniversalClient,
	scheduler domain.ScheduleReader,
) (http.Handler, error) {
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

	var authHandler http.Handler = newAuthHandler(cfg.Auth, repo, rateLimiter)
	authHandler = sloghttp.Recovery(authHandler)
	authHandler = sloghttp.NewWithConfig(slog.Default(), sloghttp.Config{
		WithRequestID: false,
		WithTraceID:   true,
		WithSpanID:    true,
	})(authHandler)
	authHandler = otelhttp.NewHandler(authHandler, "auth")
	mux.Handle("/auth/", authHandler)

	mux.Handle(contestantv1connect.NewViewerServiceHandler(
		newViewerServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewProfileServiceHandler(
		newProfileServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewContestServiceHandler(
		newContestServiceHandler(repo, scheduler),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(contestantv1connect.NewProblemServiceHandler(
		newProblemServiceHandler(repo, scheduler),
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
	mux.Handle(contestantv1connect.NewRankingServiceHandler(
		newRankingServiceHandler(repo),
		connect.WithInterceptors(interceptors...),
	))

	handler := http.Handler(mux)
	handler = session.NewHandler(sessionStore)(handler)
	handler = http.NewCrossOriginProtection().Handler(handler)
	handler = h2c.NewHandler(handler, &http2.Server{})

	return handler, nil
}
