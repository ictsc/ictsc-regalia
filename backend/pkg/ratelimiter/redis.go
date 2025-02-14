package ratelimiter

import (
	"context"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

type RateLimiter interface {
	Check(ctx context.Context, key string, window time.Duration, maxAttempts int) (bool, error)
}

type RedisRateLimiter struct {
	rdb    redis.UniversalClient
	prefix string
}

func NewRedisRateLimiter(rdb redis.UniversalClient, prefix string) *RedisRateLimiter {
	if prefix == "" {
		prefix = "rate_limit:"
	}
	return &RedisRateLimiter{
		rdb:    rdb,
		prefix: prefix,
	}
}

func (r *RedisRateLimiter) Check(ctx context.Context, key string, window time.Duration, maxAttempts int) (bool, error) {
	key = r.prefix + key
	now := time.Now().Unix()
	windowSeconds := window.Seconds()

	pipe := r.rdb.Pipeline()

	// 有効期限切れのエントリを削除
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-int64(windowSeconds), 10))
	// 現在の試行を記録
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	// 現在の試行回数を取得
	countCmd := pipe.ZCard(ctx, key)
	// 有効期限を設定
	pipe.Expire(ctx, key, window)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, errors.Wrap(err, "failed to execute redis pipeline")
	}

	attempts := countCmd.Val()
	return attempts <= int64(maxAttempts), nil
}
