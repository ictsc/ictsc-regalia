package redistest

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func SetupRedis(tb testing.TB) *redis.Client {
	tb.Helper()

	pool, err := getRedisPool()
	if err != nil {
		tb.Fatalf("failed to get redis pool: %v", err)
	}

	rdb, err := pool.acquire(tb.Context())
	if err != nil {
		tb.Fatalf("failed to acquire redis client: %v", err)
	}
	tb.Cleanup(func() {
		//nolint:usetesting // tb.Context はクリーンアップ前に終了するため使えない
		ctx, cancel := context.WithTimeout(context.Background(), redisFlushTimeout)
		defer cancel()

		if err := pool.release(ctx, rdb); err != nil {
			tb.Logf("failed to release redis client: %v", err)
		}
	})

	return rdb
}

const (
	redisImage          = "redis:7"
	redisStartupTimeout = 30 * time.Second
	redisFlushTimeout   = 10 * time.Second
	redisDatabaseSize   = 16
)

type redisPool struct {
	available chan *redis.Client
	container *tcredis.RedisContainer
}

var getRedisPool = sync.OnceValues(func() (*redisPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisStartupTimeout)
	defer cancel()

	return startRedisPool(ctx)
})

func startRedisPool(ctx context.Context) (*redisPool, error) {
	ctr, err := tcredis.Run(ctx, redisImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start redis container")
	}

	redisURL, err := ctr.ConnectionString(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get redis connection string")
	}
	baseCfg, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse redis connection string")
	}

	available := make(chan *redis.Client, redisDatabaseSize)

	var wg sync.WaitGroup
	errs := make([]error, redisDatabaseSize)
	for i := range redisDatabaseSize {
		wg.Add(1)
		go func() {
			defer wg.Done()

			cfg := *baseCfg
			cfg.DB = i
			rdb := redis.NewClient(&cfg)

			if err := rdb.Ping(ctx).Err(); err != nil {
				errs[i] = errors.Wrapf(err, "failed to ping redis db %d", i)
				return
			}

			available <- rdb
		}()
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &redisPool{
		available: available,
		container: ctr,
	}, nil
}

func (p *redisPool) acquire(ctx context.Context) (*redis.Client, error) {
	select {
	case rdb := <-p.available:
		return rdb, nil
	case <-ctx.Done():
		return nil, ctx.Err() //nolint:wrapcheck
	}
}

func (p *redisPool) release(ctx context.Context, rdb *redis.Client) error {
	if err := rdb.FlushDB(ctx).Err(); err != nil {
		return errors.Wrap(err, "failed to flush redis")
	}
	p.available <- rdb
	return nil
}
