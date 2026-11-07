package httpserver

import (
	"context"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

func (h *Handler) subscribeDeploymentEvents(ctx context.Context) (service.Subscription, error) {
	return subscribe(ctx, h.service.Events)
}

func closeDeploymentSubscription(subscription service.Subscription) {
	closeSubscription(subscription)
}
