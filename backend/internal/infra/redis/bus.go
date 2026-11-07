package redis

import (
	"context"

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

var _ service.EventBus = (*EventBus)(nil)
