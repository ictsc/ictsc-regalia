package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	redislib "github.com/redis/go-redis/v9"
)

var consumeSessionScript = redislib.NewScript(`
local raw = redis.call('GET', KEYS[1])
if not raw then
  return nil
end
local ok, value = pcall(cjson.decode, raw)
if not ok then
  return redis.error_reply('invalid session payload')
end
if value.kind ~= ARGV[1] then
  return nil
end
redis.call('DEL', KEYS[1])
return raw
`)

func (s *Store) Consume(ctx context.Context, token string, expected session.Kind) (session.Data, error) {
	if token == "" {
		return session.Data{}, session.ErrNotFound
	}
	value, err := consumeSessionScript.Run(ctx, s.client, []string{sessionPrefix + token}, string(expected)).Result()
	if err == redislib.Nil {
		return session.Data{}, session.ErrNotFound
	}
	if err != nil {
		return session.Data{}, fmt.Errorf("consume session: %w", err)
	}
	encoded, ok := value.(string)
	if !ok {
		return session.Data{}, fmt.Errorf("consume session: unexpected Redis result %T", value)
	}
	var data session.Data
	if err := json.Unmarshal([]byte(encoded), &data); err != nil {
		return session.Data{}, fmt.Errorf("decode consumed session: %w", err)
	}
	return data, nil
}
