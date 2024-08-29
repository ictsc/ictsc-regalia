package errors

import (
	"context"

	"connectrpc.com/connect"
)

// NewErrorInterceptor Connect用エラーハンドラ
func NewErrorInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			res, err := next(ctx, req)

			if err == nil {
				return res, nil
			}

			switch {
			case IsFlag(err, ErrBadArgument):
				return res, connect.NewError(connect.CodeInvalidArgument, err)
			case IsFlag(err, ErrNotFound):
				return res, connect.NewError(connect.CodeNotFound, err)
			case IsFlag(err, ErrAlreadyExists):
				return res, connect.NewError(connect.CodeAlreadyExists, err)
			default:
				return res, connect.NewError(connect.CodeUnknown, err)
			}
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
