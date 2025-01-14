package scoreserver

import (
	"context"
	"log/slog"
	"net/http"
	"net/netip"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgxutil"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant"
	"github.com/redis/go-redis/v9"
)

type ScoreServer struct {
	adminServer      *http.Server
	contestantServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*ScoreServer, error) {
	db := pgxutil.NewDBx(cfg.PgConfig, pgxutil.WithOTel(true))
	rdb := redis.NewClient(&cfg.Redis)

	adminHandler, err := admin.New(ctx, cfg.AdminAPI, db)
	if err != nil {
		return nil, err
	}

	contestantHandler, err := contestant.New(ctx, cfg.ContestantAPI, rdb)
	if err != nil {
		return nil, err
	}

	return &ScoreServer{
		adminServer:      newServer(cfg.AdminAPI.Address, adminHandler),
		contestantServer: newServer(cfg.ContestantAPI.Address, contestantHandler),
	}, nil
}

func newServer(addr netip.AddrPort, handler http.Handler) *http.Server {
	//nolint:mnd // ここに定数を閉じ込める
	return &http.Server{
		Addr:              addr.String(),
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024,
	}
}

func (s *ScoreServer) Start(ctx context.Context) error {
	go func() {
		slog.InfoContext(ctx, "Starting admin API", "address", s.adminServer.Addr)
		if err := s.adminServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "Failed to start admin API", "error", err)
		}
	}()

	go func() {
		slog.InfoContext(ctx, "Starting contestant API", "address", s.contestantServer.Addr)
		if err := s.contestantServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "Failed to start contestant API", "error", err)
		}
	}()

	return nil
}

func (s *ScoreServer) Shutdown(ctx context.Context) error {
	shutdownFns := []func(context.Context) error{
		func(ctx context.Context) error {
			slog.DebugContext(ctx, "Shutting down admin API")
			if err := s.adminServer.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "failed to shutdown admin API")
			}
			return nil
		},
		func(ctx context.Context) error {
			slog.DebugContext(ctx, "Shutting down contestant API")
			if err := s.contestantServer.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "failed to shutdown contestant API")
			}
			return nil
		},
	}
	errs := make([]error, len(shutdownFns))
	var wg sync.WaitGroup
	for i, fn := range shutdownFns {
		wg.Add(1)
		go func() {
			errs[i] = fn(ctx)
			wg.Done()
		}()
	}
	wg.Wait()

	return errors.Join(errs...)
}
