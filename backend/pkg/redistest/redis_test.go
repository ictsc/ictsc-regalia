package redistest_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/redistest"
	"github.com/redis/go-redis/v9"
)

func TestSimple(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	rdb := redistest.SetupRedis(t)

	if err := rdb.Set(ctx, "key", "value", 0).Err(); err != nil {
		t.Fatalf("failed to set: %v", err)
	}

	result := rdb.Get(ctx, "key")
	if err := result.Err(); err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	if val := result.Val(); val != "value" {
		t.Errorf("unexpected value: %v", val)
	}
}

func TestParalelUse(t *testing.T) {
	t.Parallel()

	const NumTest = 4

	var (
		count int
		cond  = sync.NewCond(&sync.Mutex{})
	)
	waitForOthers := func() {
		cond.L.Lock()
		count++
		if count < NumTest {
			cond.Wait()
		} else {
			count = 0
			cond.Broadcast()
		}
		cond.L.Unlock()
	}

	// RedisのDBを並列で利用しても問題ないことを確認する
	for i := range NumTest {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			rdb := redistest.SetupRedis(t)

			waitForOthers()

			if err := rdb.Set(ctx, "key", i, 0).Err(); err != nil {
				t.Fatalf("failed to set: %v", err)
			}

			waitForOthers()

			result := rdb.Get(ctx, "key")
			if err := result.Err(); err != nil {
				t.Fatalf("failed to get: %v", err)
			}
			if val := result.Val(); val != strconv.Itoa(i) {
				t.Errorf("unexpected value: %v", val)
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	t.Parallel()

	// RedisのDBを使い切り，使い回されたときに環境が空になっていることを確認する
	// DBプールは16個のDBを持つのでそれ以上のテストを行う
	for i := range 8 {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			var rdb *redis.Client
			for range 4 {
				rdb = redistest.SetupRedis(t)
			}

			if err := rdb.Get(ctx, "key").Err(); !errors.Is(err, redis.Nil) {
				t.Errorf("expect: %v, got: %v", redis.Nil, err)
			}

			if err := rdb.Set(ctx, "key", i, 0).Err(); err != nil {
				t.Errorf("failed to set: %v", err)
			}
		})
	}
}
