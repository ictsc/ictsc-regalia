package client_test

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-outlands/backend/pkg/connect/client"
	"github.com/stretchr/testify/assert"
)

func TestRetry_Do(t *testing.T) {
	t.Parallel()

	prodNewFunc := func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) struct{} {
		assert.Equal(t, http.DefaultClient, httpClient)

		return struct{}{}
	}
	devNewFunc := func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) struct{} {
		assert.NotEqual(t, http.DefaultClient, httpClient)

		return struct{}{}
	}

	type args struct {
		ctx      context.Context
		callback func(ctx context.Context, client struct{}) error
	}

	tests := []struct {
		name      string
		r         *client.Retry[struct{}]
		args      args
		setup     func(args *args, counter *int)
		wantCount int
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			r:    client.NewRetry(prodNewFunc, "http://example.com", 0, false),
			args: args{
				ctx:      context.Background(),
				callback: nil,
			},
			setup: func(args *args, counter *int) {
				args.callback = func(ctx context.Context, client struct{}) error {
					*counter++

					return nil
				}
			},
			wantCount: 1,
			assertion: assert.NoError,
		},
		{
			name: "success (dev)",
			r:    client.NewRetry(devNewFunc, "http://example.com", 0, true),
			args: args{
				ctx:      context.Background(),
				callback: nil,
			},
			setup: func(args *args, counter *int) {
				args.callback = func(ctx context.Context, client struct{}) error {
					*counter++

					return nil
				}
			},
			wantCount: 1,
			assertion: assert.NoError,
		},
		{
			name: "return error",
			r:    client.NewRetry(prodNewFunc, "http://example.com", 0, false),
			args: args{
				ctx:      context.Background(),
				callback: nil,
			},
			setup: func(args *args, counter *int) {
				args.callback = func(ctx context.Context, client struct{}) error {
					*counter++

					return connect.NewError(connect.CodeInternal, nil)
				}
			},
			wantCount: 1,
			assertion: assert.Error,
		},
		{
			name: "retry 10 times",
			r:    client.NewRetry(prodNewFunc, "http://example.com", 0, false),
			args: args{
				ctx:      context.Background(),
				callback: nil,
			},
			setup: func(args *args, counter *int) {
				args.callback = func(ctx context.Context, client struct{}) error {
					*counter++

					return connect.NewError(connect.CodeUnavailable, nil)
				}
			},
			wantCount: 10,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			counter := 0
			tt.setup(&tt.args, &counter)

			tt.assertion(t, tt.r.Do(tt.args.ctx, tt.args.callback))
			assert.Equal(t, tt.wantCount, counter)
		})
	}
}
