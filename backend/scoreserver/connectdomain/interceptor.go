package connectdomain

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
)

func NewErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			res, err := next(ctx, req)
			if err != nil {
				return res, connectError(err)
			}
			return res, nil
		}
	}
}

func NewLoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			procedure := req.Spec().Procedure
			peer := req.Peer()

			logger := slog.Default().With(
				"procedure", procedure,
				"protocol", peer.Protocol,
				"peer", peer.Addr,
			)

			logger.InfoContext(ctx, "Call started")

			res, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				lvl := slog.LevelError
				if cErr := new(connect.Error); errors.As(err, &cErr) {
					lvl = serverCodeToLevel(cErr.Code())
				}
				if sErr := new(sanitizedError); errors.As(err, &sErr) {
					err = sErr.err
				}

				logger.Log(ctx, lvl, "Call finished",
					"error", err,
					"duration_ms", duration.Milliseconds(),
				)
				return res, err
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
		connect.CodeInvalidArgument, connect.CodeUnauthenticated, connect.CodePermissionDenied,
		connect.CodeFailedPrecondition:
		return slog.LevelInfo
	case connect.CodeDeadlineExceeded, connect.CodeResourceExhausted, connect.CodeAborted,
		connect.CodeOutOfRange, connect.CodeUnavailable:
		return slog.LevelWarn
	case connect.CodeUnknown, connect.CodeUnimplemented,
		connect.CodeInternal, connect.CodeDataLoss:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
