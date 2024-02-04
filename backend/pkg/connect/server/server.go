// Package server Connectサーバーユーティリティ
package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type serverType int

const (
	// TypeExternal 外部向けサーバー
	TypeExternal serverType = iota
	// TypeInternal 内部向けサーバー
	TypeInternal
)

type registerer struct {
	mux             *http.ServeMux
	commonOpt       []connect.HandlerOption
	authInterceptor connect.HandlerOption
}

// New Connectサーバーを作成する
func New(
	dev bool,
	srvType serverType,
	addr string,
	register func(reg *registerer),
) (*http.Server, func()) {
	mux := http.NewServeMux()
	reg := &registerer{
		mux:             mux,
		commonOpt:       []connect.HandlerOption{},
		authInterceptor: nil,
	}

	register(reg)

	log.Print(dev, srvType)

	srv := &http.Server{ // nolint:exhaustruct
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: time.Second,
	}

	// 開発環境の場合はh2c(TLS無しのHTTP/2)を有効にする
	if dev {
		srv.Handler = h2c.NewHandler(mux, &http2.Server{}) // nolint:exhaustruct
	}

	return srv, func() {
		_ = srv.Shutdown(context.Background())
	}
}

// Register Connectサービスを認証無しで登録する
func Register[H any](
	reg *registerer,
	newServiceHandler func(svc H, opts ...connect.HandlerOption) (string, http.Handler),
	svc H,
	opts ...connect.HandlerOption,
) {
	opts = append(reg.commonOpt, opts...)
	reg.mux.Handle(newServiceHandler(svc, opts...))
}

// RegisterWithAuth Connectサービスを認証付きで登録する
func RegisterWithAuth[H any](
	reg *registerer,
	newServiceHandler func(svc H, opts ...connect.HandlerOption) (string, http.Handler),
	svc H,
	opts ...connect.HandlerOption,
) {
	opts = append(append(reg.commonOpt, reg.authInterceptor), opts...)
	reg.mux.Handle(newServiceHandler(svc, opts...))
}
