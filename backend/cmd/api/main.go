package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

func health(
	ctx context.Context,
	input *struct{},
) (*HealthOutput, error) {
	output := &HealthOutput{}
	output.Body.Status = "ok"

	return output, nil
}

func main() {
	// chiルータの初期化
	router := chi.NewRouter()

	// humaのAPIインスタンスを作成
	api := humachi.New(
		router,
		huma.DefaultConfig(
			"ICTSC 2026 Score Server API",
			"0.0.1",
		),
	)

	// ヘルスチェックエンドポイントの定義
	huma.Get(api, "/health", health)

	// サーバの設定と起動
	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
