package config

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultAddress           = ":8080"
	defaultReadHeaderTimeout = 5 * time.Second
)

// ConfigはAPIサーバーの設定を表す
type Config struct {
	Address           string
	ReadHeaderTimeout time.Duration
}

// Loadは環境変数から設定を読み込む
func Load() (Config, error) {
	config := Config{
		Address:           envOrDefault("ICTSC_API_ADDRESS", defaultAddress),
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	if value := os.Getenv("ICTSC_API_READ_HEADER_TIMEOUT"); value != "" {
		timeout, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("ICTSC_API_READ_HEADER_TIMEOUTの値が不正です: %w", err)
		}
		if timeout <= 0 {
			return Config{}, fmt.Errorf("ICTSC_API_READ_HEADER_TIMEOUTは0より大きい値にしてください")
		}
		config.ReadHeaderTimeout = timeout
	}

	return config, nil
}

// envOrDefaultは環境変数が空の場合にデフォルト値を返す
func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
