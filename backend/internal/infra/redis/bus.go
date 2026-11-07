package redis

import (
	"context"

	redislib "github.com/redis/go-redis/v9"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	deploymentChannel = "ictsc:deployment-events"
	answerChannel     = "ictsc:answer-events"
)

// EventBus carries raw events over one Redis Pub/Sub channel so that every
// API instance sees them.
type EventBus struct {
	store   *Store
	channel string
}

func NewDeploymentBus(store *Store) *EventBus {
	return &EventBus{store: store, channel: deploymentChannel}
}

func NewAnswerBus(store *Store) *EventBus { return &EventBus{store: store, channel: answerChannel} }

func (b *EventBus) Publish(ctx context.Context, event []byte) error {
	return b.store.Publish(ctx, b.channel, event)
}

func (b *EventBus) Subscribe(ctx context.Context) (service.Subscription, error) {
	inner, err := b.store.Subscribe(ctx, b.channel)
	if err != nil {
		return nil, err
	}
	subscription := &busSubscription{inner: inner, messages: make(chan []byte, 64)}
	go subscription.forward(ctx)
	return subscription, nil
}

type busSubscription struct {
	inner    *Subscription
	messages chan []byte
}

func (s *busSubscription) forward(ctx context.Context) {
	defer close(s.messages)
	for {
		select {
		case <-ctx.Done():
			return
		case received, ok := <-s.inner.C:
			if !ok {
				return
			}
			var payload []byte
			switch value := received.(type) {
			case *redislib.Message:
				payload = []byte(value.Payload)
			case *redislib.Subscription:
				// go-redis reconnected and dropped what was published in
				// between. Ending the subscription ends the SSE response, so
				// the browser reconnects and refetches from the database.
				return
			default:
				continue
			}
			select {
			case s.messages <- payload:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (s *busSubscription) Messages() <-chan []byte { return s.messages }
func (s *busSubscription) Close() error            { return s.inner.Close() }

var _ service.EventBus = (*EventBus)(nil)
