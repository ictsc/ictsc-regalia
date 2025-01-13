package admin

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(ctx context.Context, cfg config.AdminAPI, db *sqlx.DB) (http.Handler, error) {
	interceptors := []connect.Interceptor{
		connectutil.NewOtelInterceptor(),
		connectutil.NewSlogInterceptor(),
	}

	mux := http.NewServeMux()

	mux.Handle(adminv1connect.NewTeamServiceHandler(
		NewTeamServiceHandler(db),
		connect.WithInterceptors(interceptors...),
	))

	checker := grpchealth.NewStaticChecker("admin.v1.TeamService")
	mux.Handle(grpchealth.NewHandler(checker))

	handler := http.Handler(mux)

	authenticator, err := auth.NewJWTAuthenticator(ctx, cfg.Authn)
	if err != nil {
		return nil, err
	}
	handler = auth.WithAuthn(handler, authenticator)

	// gRPC requires HTTP/2
	handler = h2c.NewHandler(handler, &http2.Server{})

	return handler, nil
}
