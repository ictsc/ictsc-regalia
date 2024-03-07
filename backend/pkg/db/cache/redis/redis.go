package redis

import (
	"strconv"

	"github.com/ictsc/ictsc-outlands/backend/pkg/db/cache"
	"github.com/redis/go-redis/v9"
)

// Redis Redisクライアント
type Redis[V any] struct {
	c   *redis.Client
	srv string
}

var _ cache.DB[struct{}] = (*Redis[struct{}])(nil)

// Config Redis接続設定
type Config struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database int

	Service string // Redisを利用しているサービス名
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
		srv: conf.Service,
	}
}
