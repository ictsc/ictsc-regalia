package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// HealthOutputはヘルスチェックのレスポンスを表す
type HealthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

// registerHealthRoutesはヘルスチェックのエンドポイントを登録する
func registerHealthRoutes(api huma.API) {
	huma.Get(api, "/health", health)
}

// healthはサーバーが稼働していることを返す
func health(
	ctx context.Context,
	input *struct{},
) (*HealthOutput, error) {
	output := &HealthOutput{}
	output.Body.Status = "ok"

	return output, nil
}
