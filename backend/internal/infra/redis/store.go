package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	redislib "github.com/redis/go-redis/v9"
)

const sessionPrefix = "ictsc:session:"

type Store struct {
	client *redislib.Client
}

func Open(ctx context.Context, rawURL string) (*Store, error) {
	options, err := redislib.ParseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}
	store := &Store{client: redislib.NewClient(options)}
	if err := store.client.Ping(ctx).Err(); err != nil {
		_ = store.client.Close()
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	return store, nil
}

func New(client *redislib.Client) *Store { return &Store{client: client} }

func (s *Store) Close() error { return s.client.Close() }

func (s *Store) Create(ctx context.Context, data session.Data, ttl time.Duration) (string, error) {
	token, err := session.Token()
	if err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	data.CreatedAt = time.Now().UTC()
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("encode session: %w", err)
	}
	if err := s.client.Set(ctx, sessionPrefix+token, encoded, ttl).Err(); err != nil {
		return "", fmt.Errorf("save session: %w", err)
	}
	return token, nil
}

func (s *Store) Get(ctx context.Context, token string, expected session.Kind) (session.Data, error) {
	encoded, err := s.client.Get(ctx, sessionPrefix+token).Bytes()
	if err == redislib.Nil {
		return session.Data{}, session.ErrNotFound
	}
	if err != nil {
		return session.Data{}, fmt.Errorf("read session: %w", err)
	}
	var data session.Data
	if err := json.Unmarshal(encoded, &data); err != nil {
		return session.Data{}, fmt.Errorf("decode session: %w", err)
	}
	if data.Kind != expected {
		return session.Data{}, session.ErrNotFound
	}
	return data, nil
}

func (s *Store) Delete(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.client.Del(ctx, sessionPrefix+token).Err()
}

func (s *Store) Publish(ctx context.Context, channel string, event []byte) error {
	return s.client.Publish(ctx, channel, event).Err()
}

type Subscription struct {
	pubsub *redislib.PubSub
	C      <-chan *redislib.Message
}

func (s *Store) Subscribe(ctx context.Context, channel string) (*Subscription, error) {
	pubsub := s.client.Subscribe(ctx, channel)
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return nil, err
	}
	return &Subscription{pubsub: pubsub, C: pubsub.Channel(redislib.WithChannelSize(64))}, nil
}

func (s *Subscription) Close() error { return s.pubsub.Close() }

var _ session.Store = (*Store)(nil)
