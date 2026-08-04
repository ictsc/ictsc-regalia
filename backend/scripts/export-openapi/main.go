package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ictsc/ictsc-regalia/backend/internal/handler"
)

func main() {
	// 出力先を指定できるようにしてOpenAPIを生成する
	outputPath := flag.String("output", "openapi.json", "出力先のファイルパス")
	flag.Parse()

	app := handler.New()
	data, err := json.MarshalIndent(app.API.OpenAPI(), "", "  ")
	if err != nil {
		panic(fmt.Errorf("OpenAPIのJSON化に失敗しました: %w", err))
	}

	data = append(data, '\n')
	if err := os.WriteFile(*outputPath, data, 0o644); err != nil {
		panic(fmt.Errorf("OpenAPIの書き込みに失敗しました: %w", err))
	}

	fmt.Printf("OpenAPIを%sに出力しました\n", *outputPath)
}
