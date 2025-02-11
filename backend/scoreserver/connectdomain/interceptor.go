package connectdomain

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
)

func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			procedure := req.Spec().Procedure
			peer := req.Peer()

			span := trace.SpanFromContext(ctx)
			logger := slog.Default().With(
				"procedure", procedure,
				"protocol", peer.Protocol,
				"peer", peer.Addr,
				"trace_id", span.SpanContext().TraceID().String(),
				"span_id", span.SpanContext().SpanID().String(),
			)

			logger.InfoContext(ctx, "Call started")

			res, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				connectErr := connectError(err)
				lvl := serverCodeToLevel(connectErr.Code())
				logger.Log(ctx, lvl, "Call finished",
					"error", err,
					"duration_ms", duration.Milliseconds(),
				)
				return res, connectErr
			}

			logger.InfoContext(ctx, "Call finished",
				"duration_ms", duration.Milliseconds(),
			)
			return res, nil
		}
	}
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
