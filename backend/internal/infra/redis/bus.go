package redis

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const deploymentChannel = "ictsc:deployment-events"

type DeploymentBus struct{ store *Store }

func NewDeploymentBus(store *Store) *DeploymentBus { return &DeploymentBus{store: store} }

func (b *DeploymentBus) Publish(ctx context.Context, event []byte) error {
	return b.store.Publish(ctx, deploymentChannel, event)
}

func (b *DeploymentBus) Subscribe(ctx context.Context) (service.Subscription, error) {
	inner, err := b.store.Subscribe(ctx, deploymentChannel)
	if err != nil {
		return nil, err
	}
	subscription := &deploymentSubscription{inner: inner, messages: make(chan []byte, 64)}
	go subscription.forward(ctx)
	return subscription, nil
}

type deploymentSubscription struct {
	inner    *Subscription
	messages chan []byte
}

func (s *deploymentSubscription) forward(ctx context.Context) {
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

func (s *deploymentSubscription) Messages() <-chan []byte { return s.messages }
func (s *deploymentSubscription) Close() error            { return s.inner.Close() }

var _ service.EventBus = (*DeploymentBus)(nil)
