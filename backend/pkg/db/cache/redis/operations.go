package redis

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// Get 指定したキーの値を取得する
func (r *Redis[V]) Get(ctx context.Context, key string) (*V, error) { // nolint:ireturn
	valueStr, err := r.c.Get(ctx, r.srv+"-"+key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	} else if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	var value V

	err = sonic.Unmarshal([]byte(valueStr), &value)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return &value, nil
}

// Set 指定したキーに値を設定する
//
//	ttlはミリ秒単位で指定
func (r *Redis[V]) Set(ctx context.Context, key string, value V, ttl int) error {
	valueStr, err := sonic.Marshal(value)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	err = r.c.Set(ctx, r.srv+"-"+key, string(valueStr), time.Millisecond*time.Duration(ttl)).Err()
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	return nil
}
