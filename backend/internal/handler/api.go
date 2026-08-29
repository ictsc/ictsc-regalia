package handler

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/ictsc/ictsc-regalia/backend/internal/handler/admin"
)

// AppはAPIとHTTPルーターをまとめて保持する
type App struct {
	Router *chi.Mux
	API    huma.API
}

// NewはHuma APIを作成して各エンドポイントを登録する
func New() *App {
	router := chi.NewRouter()
	api := humachi.New(
		router,
		huma.DefaultConfig(
			"ICTSC 2026 Score Server API",
			"0.0.1",
		),
	)

	registerHealthRoutes(api)
	admin.RegisterTeamRoutes(api)

	return &App{
		Router: router,
		API:    api,
	}
}
