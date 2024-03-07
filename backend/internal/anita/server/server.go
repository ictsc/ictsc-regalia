package server

import (
	"net/http"

	"github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1/v1connect"
	"github.com/ictsc/ictsc-outlands/backend/pkg/connect/server"
)

// Config サーバー設定
type Config struct {
	Dev  bool
	Addr string
}

// NewServer サーバーを作成する
func NewServer(conf *Config, user v1connect.UserServiceHandler, team v1connect.TeamServiceHandler) (*http.Server, func()) {
	return server.New(conf.Dev, server.TypeInternal, conf.Addr, func(reg *server.Registerer) {
		server.Register(reg, v1connect.NewUserServiceHandler, user, nil)
		server.Register(reg, v1connect.NewTeamServiceHandler, team, nil)
	})
}
