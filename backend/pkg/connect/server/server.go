package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/ictsc/ictsc-outlands/backend/pkg/log"
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

// Registerer Connectサーバーハンドラー登録構造体
type Registerer struct {
	mux             *http.ServeMux
	commonOpt       []connect.HandlerOption
	authInterceptor connect.HandlerOption
}

// New Connectサーバーを作成する
func New(
	dev bool,
	_ serverType,
	port int,
	register func(reg *Registerer),
) (*http.Server, func()) {
	mux := http.NewServeMux()
	reg := &Registerer{
		mux: mux,
		commonOpt: []connect.HandlerOption{
			connect.WithInterceptors(
				log.NewLoggerInterceptor(log.NewLogger(dev)),
				errors.NewErrorInterceptor(),
			),
		},
		authInterceptor: nil,
	}

	register(reg)

	// Healthcheckのため、pingハンドラを追加
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}))

	srv := &http.Server{ // nolint:exhaustruct
		Addr:        "0.0.0.0:" + strconv.Itoa(port),
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
	reg *Registerer,
	newServiceHandler func(svc H, opts ...connect.HandlerOption) (string, http.Handler),
	svc H,
	opts ...connect.HandlerOption,
) {
	opts = append(reg.commonOpt, opts...)
	reg.mux.Handle(newServiceHandler(svc, opts...))
}

// RegisterWithAuth Connectサービスを認証付きで登録する
func RegisterWithAuth[H any](
	reg *Registerer,
	newServiceHandler func(svc H, opts ...connect.HandlerOption) (string, http.Handler),
	svc H,
	opts ...connect.HandlerOption,
) {
	opts = append(append(reg.commonOpt, reg.authInterceptor), opts...)
	reg.mux.Handle(newServiceHandler(svc, opts...))
}
