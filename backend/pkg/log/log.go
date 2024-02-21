// Package log ログユーティリティー
package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Config ロガー設定
type Config struct {
	Dev bool
}

// NewLogger ロガー生成
func NewLogger(conf *Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if conf.Dev {
		// nolint:exhaustruct
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger
}
