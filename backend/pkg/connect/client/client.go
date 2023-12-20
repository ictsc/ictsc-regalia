// Package client Connectクライアントラッパー実装
package client

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
)

const defaultRetry = 10

// Retry リトライ機能付きConnectクライアント
type Retry[C any] struct {
	client C
	retry  int
}

// NewRetry Retryを生成
func NewRetry[C any](
	newFunc func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) C,
	baseURL string,
	retry int,
	dev bool,
) *Retry[C] {
	client := http.DefaultClient
	// 開発環境の場合はh2c(TLS無しのHTTP/2)を有効にする
	if dev {
		client = &http.Client{ // nolint:exhaustruct
			Transport: &http2.Transport{ // nolint:exhaustruct
				AllowHTTP: true,
				DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr) // nolint:wrapcheck
				},
			},
		}
	}

	if retry == 0 {
		retry = defaultRetry
	}

	return &Retry[C]{
		client: newFunc(client, baseURL, connect.WithSendGzip()),
		retry:  retry,
	}
}

// Do リトライ付きでcallbackを実行
func (c *Retry[C]) Do(ctx context.Context, callback func(ctx context.Context, client C) error) error {
	var err error

	for i := 0; i < c.retry; i++ {
		err = callback(ctx, c.client)

		switch {
		case err == nil:
			return nil
		case connect.CodeOf(err) == connect.CodeUnavailable:
			continue
		default:
			return err
		}
	}

	return err
}
