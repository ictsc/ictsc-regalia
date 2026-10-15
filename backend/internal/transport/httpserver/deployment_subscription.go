package httpserver

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

func (h *Handler) subscribeDeploymentEvents(ctx context.Context) (service.Subscription, error) {
	if h.service.Events == nil {
		return nil, nil
	}
	return h.service.Events.Subscribe(ctx)
}

func closeDeploymentSubscription(subscription service.Subscription) {
	if subscription != nil {
		_ = subscription.Close()
	}
}
