package contestant

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/auth"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(ctx context.Context, cfg config.ContestantAPI, rdb redis.UniversalClient) (http.Handler, error) {
	sessStore, err := redisstore.NewRedisStore(ctx, rdb)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create redis store")
	}

	mux := http.NewServeMux()

	handler := http.Handler(mux)

	mux.Handle("/auth/", auth.NewHandler(cfg.Auth, sessStore))

	handler = h2c.NewHandler(handler, &http2.Server{})

	return handler, nil
}
