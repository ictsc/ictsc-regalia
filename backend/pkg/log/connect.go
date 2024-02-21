package log

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/rs/zerolog"
)

// NewLoggerInterceptor Connect用ロガー
func NewLoggerInterceptor(logger zerolog.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			start := time.Now()

			res, err := next(ctx, req)

			if err != nil {
				logger.
					Error().
					Str("procedure", req.Spec().Procedure).
					Str("latency", time.Since(start).String()).
					Msg(errors.Sprint(err))
			} else {
				logger.
					Info().
					Str("procedure", req.Spec().Procedure).
					Str("latency", time.Since(start).String()).
					Msg("")
			}

			return res, err
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
