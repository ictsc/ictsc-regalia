package redis

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// Get 指定したキーの値を取得する
func (r *Redis[K, V]) Get(ctx context.Context, key K) (*V, error) { // nolint:ireturn
	keyStr, err := sonic.Marshal(key)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	valueStr, err := r.c.Get(ctx, string(keyStr)).Result()
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
func (r *Redis[K, V]) Set(ctx context.Context, key K, value V, ttl int) error {
	keyStr, err := sonic.Marshal(key)
	if err != nil {
		return errors.Wrap(err)
	}

	valueStr, err := sonic.Marshal(value)
	if err != nil {
		return errors.Wrap(err)
	}

	err = r.c.Set(ctx, string(keyStr), string(valueStr), time.Millisecond*time.Duration(ttl)).Err()

	return errors.Wrap(err)
}
