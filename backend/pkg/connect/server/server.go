// Package server Connectサーバーユーティリティ
package server

import (
	"log"
	"net/http"

	"connectrpc.com/connect"
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

// NewMux Connectサーバールーターを作成する
func NewMux(
	dev bool,
	srvType serverType,
	register func(reg *registerer),
) *http.ServeMux {
	mux := http.NewServeMux()
	reg := &registerer{
		mux:             mux,
		commonOpt:       []connect.HandlerOption{},
		authInterceptor: nil,
	}

	register(reg)

	log.Print(dev, srvType)

	return mux
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
