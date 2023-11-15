// Package redis Redisユーティリティー
package redis

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Config Redis接続用設定
type Config struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database int
}

// Redis Redisクライアント
type Redis[V any] struct {
	c *redis.Client
}

// New Redisクライアント生成
func New[V any](conf *Config) *Redis[V] {
	return &Redis[V]{
		c: redis.NewClient(&redis.Options{ // nolint:exhaustruct
			Addr:     conf.Hostname + ":" + strconv.Itoa(conf.Port),
			Username: conf.Username,
			Password: conf.Password,
			DB:       conf.Database,
		}),
	}
}
