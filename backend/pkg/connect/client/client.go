package client

import (
	"crypto/tls"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
)

// New connectのClientを生成
func New[C any]( // nolint:ireturn
	newFunc func(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) C,
	baseURL string,
	dev bool,
	opts ...connect.ClientOption,
) C {
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

	return newFunc(client, baseURL, append(opts, connect.WithSendGzip())...)
}
