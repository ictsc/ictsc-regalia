package connectslog

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New() connect.Interceptor {
	logger := slog.Default()
	unaryInterceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			service, method := splitProceduce(req.Spec().Procedure)

			start := now()
			resp, err := next(ctx, req)
			duration := now().Sub(start)

			code := status.Code(err)
			lvl := serverCodeToLevel(code)
			attrs := []slog.Attr{
				slog.String("grpc.service", service),
				slog.String("grpc.method", method),
				slog.String("grpc.duration", duration.String()),
				slog.Int("grpc.code", int(code)),
			}
			if err != nil {
				attrs = append(attrs, slog.String("grpc.error", err.Error()))
			}
			logger.LogAttrs(ctx, lvl, "Call finished", attrs...)

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(unaryInterceptor)
}

var now = time.Now

func serverCodeToLevel(code codes.Code) slog.Level {
	switch code {
	case codes.OK, codes.NotFound, codes.Canceled, codes.AlreadyExists, codes.InvalidArgument, codes.Unauthenticated:
		return slog.LevelInfo
	case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted,
		codes.OutOfRange, codes.Unavailable:
		return slog.LevelWarn
	case codes.Unknown, codes.Unimplemented, codes.Internal, codes.DataLoss:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func splitProceduce(proc string) (service, method string) {
	proc = strings.TrimLeft(proc, "/")
	service, method, _ = strings.Cut(proc, "/")
	return
}
