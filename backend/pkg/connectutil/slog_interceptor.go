package connectutil

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
)

func NewSlogInterceptor() connect.Interceptor {
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

			code := connect.CodeOf(err)
			lvl := serverCodeToLevel(code)
			attrs := []slog.Attr{
				slog.String("rpc.service", service),
				slog.String("rpc.method", method),
				slog.String("duration", duration.String()),
				slog.String("rpc.code", code.String()),
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

func serverCodeToLevel(code connect.Code) slog.Level {
	switch code {
	case 0, connect.CodeNotFound, connect.CodeCanceled, connect.CodeAlreadyExists,
		connect.CodeInvalidArgument, connect.CodeUnauthenticated:
		return slog.LevelInfo
	case connect.CodeDeadlineExceeded, connect.CodePermissionDenied, connect.CodeResourceExhausted,
		connect.CodeFailedPrecondition, connect.CodeAborted,
		connect.CodeOutOfRange, connect.CodeUnavailable:
		return slog.LevelWarn
	case connect.CodeUnknown, connect.CodeUnimplemented,
		connect.CodeInternal, connect.CodeDataLoss:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

//nolint:nonamedreturns // 名前が付いたほうが適切
func splitProceduce(proc string) (service, method string) {
	proc = strings.TrimLeft(proc, "/")
	service, method, _ = strings.Cut(proc, "/")
	return
}
