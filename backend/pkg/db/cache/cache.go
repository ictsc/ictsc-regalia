package cache

import "context"

// DB キャッシュ用データベースインターフェス
type DB[V any] interface {
	Get(ctx context.Context, key string) (*V, error)
	Set(ctx context.Context, key string, value V, ttl int) error
}
