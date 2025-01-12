package scoreserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgxutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin"
	adminauth "github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/jmoiron/sqlx"
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

func New(ctx context.Context, cfg *config.Config) (*ScoreServer, error) {
	db := pgxutil.NewDBx(cfg.PgConfig, pgxutil.WithOTel(true))

	adminServer, err := newAdminAPI(ctx, cfg.AdminAPI, db)
	if err != nil {
		return nil, err
	}

	return &ScoreServer{
		adminServer: adminServer,
	}, nil
}

func newAdminAPI(ctx context.Context, cfg config.AdminAPIConfig, db *sqlx.DB) (*http.Server, error) {
	authenticator, err := adminauth.NewJWTAuthenticator(ctx, cfg.Authn)
	if err != nil {
		return nil, err
	}
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
	handler = adminauth.WithAuthn(handler, authenticator)

	// gRPC requires HTTP/2
	handler = h2c.NewHandler(handler, &http2.Server{})

	return &http.Server{
		Addr:              cfg.Address.String(),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}, nil
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
