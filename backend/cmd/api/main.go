package main

import (
	"log"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/config"
	"github.com/ictsc/ictsc-regalia/backend/internal/handler"
)

func main() {
	// 環境変数からサーバーの設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Huma APIとHTTPルーターを作成する
	app := handler.New()

	// HTTPサーバーを起動する
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           app.Router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
