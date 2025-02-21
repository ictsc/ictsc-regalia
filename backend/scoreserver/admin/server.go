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
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/connectdomain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/growi"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(ctx context.Context, cfg config.AdminAPI, db *sqlx.DB) (http.Handler, error) {
	enforcer, err := auth.NewEnforcer(cfg.Authz)
	if err != nil {
		return nil, err
	}

	repo := pg.NewRepository(db)
	growiClient := &growi.Client{
		APIToken: cfg.Growi.Token,
		BaseURL:  cfg.Growi.BaseURL,
		Client:   otelhttp.DefaultClient,
	}

	interceptors := []connect.Interceptor{
		connectutil.NewOtelInterceptor(),
		connectdomain.NewErrorInterceptor(),
		connectdomain.NewLoggingInterceptor(),
	}

	mux := http.NewServeMux()

	mux.Handle(adminv1connect.NewTeamServiceHandler(
		NewTeamServiceHandler(enforcer, repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(adminv1connect.NewInvitationServiceHandler(
		NewInvitationServiceHandler(enforcer, repo),
		connect.WithInterceptors(interceptors...),
	))
	mux.Handle(adminv1connect.NewProblemServiceHandler(
		NewProblemServiceHandler(enforcer, repo, growiClient),
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

func enforce(ctx context.Context, enforcer *auth.Enforcer, obj, act string) error {
	viewer := auth.GetViewer(ctx)
	ok, err := enforcer.Enforce(viewer, obj, act)
	if err != nil {
		return err
	}
	if !ok {
		return connect.NewError(connect.CodePermissionDenied, nil)
	}
	return nil
}
