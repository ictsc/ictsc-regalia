package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/pkg/ratelimiter"
	"github.com/ictsc/ictsc-regalia/backend/pkg/redistest"
)

func TestRedisRateLimiter(t *testing.T) {
	t.Parallel()
	rdb := redistest.SetupRedis(t)
	limiter := ratelimiter.NewRedisRateLimiter(rdb, "test:")

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		window := 10 * time.Second
		maxAttempts := 3

		for i := range 3 {
			ok, err := limiter.Check(t.Context(), t.Name(), window, maxAttempts)
			if err != nil {
				t.Fatalf("attempt %d: %+v", i, err)
			}
			if !ok {
				t.Errorf("attempt %d: expected ok, but not ok", i)
			}
			time.Sleep(time.Second)
		}
	})

	t.Run("rate limit", func(t *testing.T) {
		t.Parallel()

		window := 10 * time.Second
		maxAttempts := 0

		ok, err := limiter.Check(t.Context(), t.Name(), window, maxAttempts)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if ok {
			t.Errorf("expected not ok, but ok")
		}
	})

	t.Run("sliding window", func(t *testing.T) {
		t.Parallel()

		window := 5 * time.Second
		maxAttempts := 2

		// time: 0
		// 最初に2回試行
		for i := range 2 {
			if ok, err := limiter.Check(t.Context(), t.Name(), window, maxAttempts); err != nil {
				t.Fatalf("attempt %d: %+v", i, err)
			} else if !ok {
				t.Fatalf("attempt %d: expected ok, but not ok", i)
			}
			time.Sleep(time.Second)
		}
		// time: 2
		// 3回目の試行は失敗する
		if ok, err := limiter.Check(t.Context(), t.Name(), window, maxAttempts); err != nil {
			t.Fatalf("%+v", err)
		} else if ok {
			t.Fatalf("expected not ok, but ok")
		}

		time.Sleep(5 * time.Second)
		// time: 7
		// 3回目の試行だけが記録に残っているので4回目は成功する
		if ok, err := limiter.Check(t.Context(), t.Name(), window, maxAttempts); err != nil {
			t.Fatalf("%+v", err)
		} else if !ok {
			t.Fatalf("expected ok, but not ok")
		}
	})
}
