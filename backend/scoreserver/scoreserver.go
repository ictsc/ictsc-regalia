package scoreserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"github.com/XSAM/otelsql"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	readHeaderTimeout = time.Second
	readTimeout       = 5 * time.Minute
	writeTimeout      = 5 * time.Minute
	maxHeaderBytes    = 8 * 1024
)

type ScoreServer struct {
	adminServer *http.Server
}

func New(cfg *Config) (*ScoreServer, error) {
	pgcfg, err := pgx.ParseConfig(cfg.DBDSN)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse DB URL")
	}

	db := sqlx.NewDb(otelsql.OpenDB(
		stdlib.GetConnector(*pgcfg,
			stdlib.OptionAfterConnect(func(_ context.Context, c *pgx.Conn) error {
				pgxuuid.Register(c.TypeMap())
				return nil
			}),
		),
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
			semconv.DBNamespace(pgcfg.Database),
		),
	), "pgx")

	adminServer := cfg.AdminAPI.new(db)

	return &ScoreServer{
		adminServer: adminServer,
	}, nil
}

func (cfg *AdminAPIConfig) new(db *sqlx.DB) *http.Server {
	var interceptors []connect.Interceptor

	interceptors = append(interceptors,
		connectutil.NewOtelInterceptor(),
		connectutil.NewSlogInterceptor(),
	)

	mux := http.NewServeMux()

	mux.Handle(adminv1connect.NewTeamServiceHandler(
		admin.NewTeamServiceHandler(db),
		connect.WithInterceptors(interceptors...),
	))

	checker := grpchealth.NewStaticChecker("admin.v1.TeamService")
	mux.Handle(grpchealth.NewHandler(checker))

	handler := http.Handler(mux)

	// gRPC requires HTTP/2
	handler = h2c.NewHandler(handler, &http2.Server{})

	return &http.Server{
		Addr:              cfg.Address.String(),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

func (s *ScoreServer) Start(_ context.Context) error {
	go func() {
		slog.Info("Starting admin API", "address", s.adminServer.Addr)
		if err := s.adminServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start admin API", "error", err)
		}
	}()

	return nil
}

func (s *ScoreServer) Shutdown(ctx context.Context) error {
	slog.DebugContext(ctx, "Shutting down admin API")
	if err := s.adminServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown admin API")
	}

	return nil
}
