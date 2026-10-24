package redis

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	deploymentChannel = "ictsc:deployment-events"
	realtimeChannel   = "ictsc:realtime-events"
)

// Bus is a Redis Pub/Sub channel shared by every API instance.
type Bus struct {
	store   *Store
	channel string
}

func NewDeploymentBus(store *Store) *Bus { return &Bus{store: store, channel: deploymentChannel} }

func NewRealtimeBus(store *Store) *Bus { return &Bus{store: store, channel: realtimeChannel} }

func (b *Bus) Publish(ctx context.Context, event []byte) error {
	return b.store.Publish(ctx, b.channel, event)
}

func (b *Bus) Subscribe(ctx context.Context) (service.Subscription, error) {
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
		case message, ok := <-s.inner.C:
			if !ok {
				return
			}
			payload := []byte(message.Payload)
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

var _ service.EventBus = (*Bus)(nil)
