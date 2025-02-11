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
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(ctx context.Context, cfg config.AdminAPI, db *sqlx.DB) (http.Handler, error) {
	enforcer, err := auth.NewEnforcer(cfg.Authz)
	if err != nil {
		return nil, err
	}

	repo := pg.NewRepository(db)

	interceptors := []connect.Interceptor{
		connectutil.NewOtelInterceptor(),
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

// connectError は domain のエラーを connect のエラーに変換する
func connectError(err error) error {
	if err == nil {
		return nil
	}
	var code connect.Code
	switch domain.ErrorType(err) {
	case domain.ErrTypeInternal:
		code = connect.CodeInternal
	case domain.ErrTypeInvalidArgument:
		code = connect.CodeInvalidArgument
	case domain.ErrTypeAlreadyExists:
		code = connect.CodeAlreadyExists
	case domain.ErrTypeNotFound:
		code = connect.CodeNotFound
	case domain.ErrTypeUnknown:
		fallthrough
	default:
		code = connect.CodeUnknown
	}
	return connect.NewError(code, err)
}

func enforce(ctx context.Context, enforcer *auth.Enforcer, obj, act string) error {
	viewer := auth.GetViewer(ctx)
	ok, err := enforcer.Enforce(viewer, obj, act)
	if err != nil {
		return connectError(err)
	}
	if !ok {
		return connect.NewError(connect.CodePermissionDenied, nil)
	}
	return nil
}
