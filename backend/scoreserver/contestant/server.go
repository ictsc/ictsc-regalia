package contestant

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(ctx context.Context, cfg config.ContestantAPI) (http.Handler, error) {
	interceptors := []connect.Interceptor{
		connectutil.NewOtelInterceptor(),
		connectutil.NewSlogInterceptor(),
	}

	mux := http.NewServeMux()

	mux.Handle(contestantv1connect.NewViewerServiceHandler(
		contestantv1connect.UnimplementedViewerServiceHandler{},
		connect.WithInterceptors(interceptors...),
	))

	handler := http.Handler(mux)
	handler = h2c.NewHandler(handler, &http2.Server{})

	return handler, nil
}
