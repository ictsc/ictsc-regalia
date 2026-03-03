package scoreserver

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgxutil"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/fake"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	readHeaderTimeout = time.Second
	readTimeout       = 5 * time.Minute
	writeTimeout      = 5 * time.Minute
	maxHeaderBytes    = 8 * 1024
)

type ScoreServer struct {
	adminServer      *http.Server
	contestantServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*ScoreServer, error) {
	db := pgxutil.NewDBx(cfg.PgConfig, pgxutil.WithOTel(true))

	rdb := redis.NewClient(&cfg.Redis)
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		slog.InfoContext(ctx, "Failed to instrument tracing for redis", "error", errors.WithStack(err))
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		slog.InfoContext(ctx, "Failed to instrument metrics for redis", "error", errors.WithStack(err))
	}

	repo := pg.NewRepository(db)

	var scheduleReader domain.ScheduleReader
	if cfg.FakeSchedule != nil {
		scheduleReader = (*fake.FakeScheduler)(cfg.FakeSchedule)
	} else {
		scheduleReader = repo
	}

	adminHandler, err := admin.New(ctx, cfg.AdminAPI, db, scheduleReader)
	if err != nil {
		return nil, err
	}
	adminServer := &http.Server{
		Addr:              cfg.AdminAPI.Address.String(),
		Handler:           adminHandler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	contestantHandler, err := contestant.New(ctx, cfg.ContestantAPI, db, rdb, scheduleReader)
	if err != nil {
		return nil, err
	}
	contestantServer := &http.Server{
		Addr:              cfg.ContestantAPI.Address.String(),
		Handler:           contestantHandler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	return &ScoreServer{
		adminServer:      adminServer,
		contestantServer: contestantServer,
	}, nil
}

func (s *ScoreServer) Start(_ context.Context) error {
	go func() {
		slog.Info("Starting admin API", "address", s.adminServer.Addr)
		if err := s.adminServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start admin API", "error", err)
		}
	}()
	go func() {
		slog.Info("Starting contestant API", "address", s.contestantServer.Addr)
		if err := s.contestantServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start contestant API", "error", err)
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
			defer wg.Done()
			errs[i] = fn(ctx)
		}()
	}
	wg.Wait()

	return errors.Join(errs...)
}
