package httpserver

import (
	"context"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) GetContestantWebPushConfig(context.Context, api.GetContestantWebPushConfigRequestObject) (api.GetContestantWebPushConfigResponseObject, error) {
	enabled, publicKey := h.service.WebPushConfig()
	response := api.WebPushConfig{Enabled: enabled}
	if enabled {
		response.VapidPublicKey = &publicKey
	}
	return api.GetContestantWebPushConfig200JSONResponse(response), nil
}

func (h *Handler) PutContestantWebPushSubscription(ctx context.Context, request api.PutContestantWebPushSubscriptionRequestObject) (api.PutContestantWebPushSubscriptionResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	if err := h.service.RegisterWebPush(ctx, principal(ctx).Session.ContestantName,
		request.Body.Endpoint, request.Body.Keys.P256dh, request.Body.Keys.Auth); err != nil {
		return nil, err
	}
	return api.PutContestantWebPushSubscription204Response{}, nil
}

func (h *Handler) DeleteContestantWebPushSubscription(ctx context.Context, request api.DeleteContestantWebPushSubscriptionRequestObject) (api.DeleteContestantWebPushSubscriptionResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	if err := h.service.DeleteWebPush(ctx, principal(ctx).Session.ContestantName, request.Body.Endpoint); err != nil {
		return nil, err
	}
	return api.DeleteContestantWebPushSubscription204Response{}, nil
}
