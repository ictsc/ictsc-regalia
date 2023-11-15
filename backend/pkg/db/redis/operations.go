package redis

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// Get 指定したキーの値を取得する
func (r *Redis[V]) Get(ctx context.Context, key string) (*V, error) { // nolint:ireturn
	valueStr, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var value V
	err = sonic.Unmarshal([]byte(valueStr), &value)

	return &value, errors.Wrap(err)
}

// Set 指定したキーに値を設定する
//
//	ttlはミリ秒単位で指定
func (r *Redis[V]) Set(ctx context.Context, key string, value V, ttl int) error {
	valueStr, err := sonic.Marshal(value)
	if err != nil {
		return errors.Wrap(err)
	}

	err = r.c.Set(ctx, key, string(valueStr), time.Millisecond*time.Duration(ttl)).Err()

	return errors.Wrap(err)
}
