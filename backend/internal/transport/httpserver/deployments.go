package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) ListContestantDeployments(ctx context.Context, request api.ListContestantDeploymentsRequestObject) (api.ListContestantDeploymentsResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	team, problem := contestant.TeamCode, string(request.ProblemCode)
	deployments, err := h.service.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		return nil, err
	}
	mapped := make([]api.ContestantDeployment, 0, len(deployments))
	for _, deployment := range deployments {
		mappedDeployment, err := h.toContestantDeployment(ctx, deployment)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, mappedDeployment)
	}
	sort.SliceStable(mapped, func(i, j int) bool { return mapped[i].Revision > mapped[j].Revision })
	return api.ListContestantDeployments200JSONResponse{Deployments: mapped}, nil
}

func (h *Handler) CreateContestantDeployment(ctx context.Context, request api.CreateContestantDeploymentRequestObject) (api.CreateContestantDeploymentResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	if h.service.Deployments == nil {
		return nil, core.NewError(http.StatusBadGateway, "upstream_unavailable", "SState is not configured")
	}
	deployment, err := h.service.QueueDeployment(ctx, contestant.TeamCode, request.ProblemCode, true)
	if err != nil {
		return nil, err
	}
	mapped, err := h.toContestantDeployment(ctx, deployment)
	if err != nil {
		return nil, err
	}
	return api.CreateContestantDeployment201JSONResponse{Deployment: mapped}, nil
}

func (h *Handler) ListAdminDeployments(ctx context.Context, request api.ListAdminDeploymentsRequestObject) (api.ListAdminDeploymentsResponseObject, error) {
	filter := core.DeploymentFilter{}
	if request.Params.TeamCode != nil {
		value := int64(*request.Params.TeamCode)
		filter.TeamCode = &value
	}
	if request.Params.ProblemCode != nil {
		value := string(*request.Params.ProblemCode)
		filter.ProblemCode = &value
	}
	deployments, err := h.service.Store.ListDeployments(ctx, filter)
	if err != nil {
		return nil, err
	}
	mapped := make([]api.AdminDeployment, 0, len(deployments))
	for _, deployment := range deployments {
		mapped = append(mapped, toAPIAdminDeployment(deployment))
	}
	sortAdminDeployments(mapped)
	return api.ListAdminDeployments200JSONResponse{Deployments: mapped}, nil
}

func (h *Handler) CreateAdminDeployment(ctx context.Context, request api.CreateAdminDeploymentRequestObject) (api.CreateAdminDeploymentResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	if h.service.Deployments == nil {
		return nil, core.NewError(http.StatusBadGateway, "upstream_unavailable", "SState is not configured")
	}
	if _, err := h.service.Store.GetTeam(ctx, request.Body.TeamCode); err != nil {
		return nil, err
	}
	deployment, err := h.service.QueueDeployment(ctx, request.Body.TeamCode, request.Body.ProblemCode, false)
	if err != nil {
		return nil, err
	}
	return api.CreateAdminDeployment201JSONResponse{Deployment: toAPIAdminDeployment(deployment)}, nil
}

func (h *Handler) SyncAdminDeployment(ctx context.Context, request api.SyncAdminDeploymentRequestObject) (api.SyncAdminDeploymentResponseObject, error) {
	if h.service.Deployments == nil {
		return nil, core.NewError(http.StatusBadGateway, "upstream_unavailable", "SState is not configured")
	}
	team, problem := int64(request.TeamCode), string(request.ProblemCode)
	deployments, err := h.service.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		return nil, err
	}
	if len(deployments) == 0 {
		return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Deployment was not found")
	}
	latest := deployments[0]
	for _, deployment := range deployments[1:] {
		if deployment.Revision > latest.Revision {
			latest = deployment
		}
	}
	event, err := h.service.Deployments.Status(ctx, team, problem)
	if err != nil {
		var notFound interface{ DeploymentStatusNotFound() bool }
		if errors.As(err, &notFound) && notFound.DeploymentStatusNotFound() {
			if latest.LatestStatus == core.DeploymentQueued {
				if failErr := h.service.FailQueuedDeployment(ctx, latest, "SState has no matching deployment request"); failErr != nil {
					return nil, failErr
				}
			}
			return nil, core.WrapError(http.StatusNotFound, "resource_not_found", "SState deployment status was not found", err)
		}
		return nil, core.WrapError(http.StatusBadGateway, "upstream_unavailable", "SState status lookup failed", err)
	}
	if event.EventID == "" {
		event.EventID = uuid.NewString()
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = h.service.Now()
	}
	if _, _, err := h.service.ApplyDeploymentEvent(ctx, team, problem, latest.Revision, event, true); err != nil {
		return nil, err
	}
	return api.SyncAdminDeployment204Response{}, nil
}

func (h *Handler) CreateAdminDeploymentEvent(ctx context.Context, request api.CreateAdminDeploymentEventRequestObject) (api.CreateAdminDeploymentEventResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	event := core.DeploymentEvent{
		EventID: request.Body.EventId.String(), OccurredAt: request.Body.OccurredAt,
		Status: core.DeploymentStatus(request.Body.Status), Message: request.Body.Message,
	}
	_, duplicate, err := h.service.ApplyDeploymentEvent(ctx, request.TeamCode, request.ProblemCode, request.Revision, event, false)
	if err != nil {
		return nil, err
	}
	response := api.DeploymentEventResponse{Event: toAPIEvent(event)}
	if duplicate {
		return api.CreateAdminDeploymentEvent200JSONResponse(response), nil
	}
	return api.CreateAdminDeploymentEvent201JSONResponse(response), nil
}

func (h *Handler) StreamContestantDeployments(ctx context.Context, request api.StreamContestantDeploymentsRequestObject) (api.StreamContestantDeploymentsResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	team, problem := contestant.TeamCode, string(request.ProblemCode)
	subscription, err := h.subscribeDeploymentEvents(ctx)
	if err != nil {
		closeDeploymentSubscription(subscription)
		return nil, err
	}
	deployments, err := h.service.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &team, ProblemCode: &problem})
	if err != nil {
		closeDeploymentSubscription(subscription)
		return nil, err
	}
	snapshot := make([]api.ContestantDeployment, 0, len(deployments))
	for _, deployment := range deployments {
		mapped, err := h.toContestantDeployment(ctx, deployment)
		if err != nil {
			closeDeploymentSubscription(subscription)
			return nil, err
		}
		snapshot = append(snapshot, mapped)
	}
	sortContestantDeployments(snapshot)
	body := h.contestantDeploymentStream(ctx, team, problem, snapshot, subscription)
	cache, buffering := "no-cache, no-store", "no"
	return api.StreamContestantDeployments200TexteventStreamResponse{
		Body: body, Headers: api.StreamContestantDeployments200ResponseHeaders{CacheControl: &cache, XAccelBuffering: &buffering},
	}, nil
}

func (h *Handler) StreamAdminDeployments(ctx context.Context, request api.StreamAdminDeploymentsRequestObject) (api.StreamAdminDeploymentsResponseObject, error) {
	filter := core.DeploymentFilter{}
	if request.Params.TeamCode != nil {
		value := int64(*request.Params.TeamCode)
		filter.TeamCode = &value
	}
	if request.Params.ProblemCode != nil {
		value := string(*request.Params.ProblemCode)
		filter.ProblemCode = &value
	}
	subscription, err := h.subscribeDeploymentEvents(ctx)
	if err != nil {
		closeDeploymentSubscription(subscription)
		return nil, err
	}
	deployments, err := h.service.Store.ListDeployments(ctx, filter)
	if err != nil {
		closeDeploymentSubscription(subscription)
		return nil, err
	}
	snapshot := make([]api.AdminDeployment, 0, len(deployments))
	for _, deployment := range deployments {
		snapshot = append(snapshot, toAPIAdminDeployment(deployment))
	}
	sortAdminDeployments(snapshot)
	body := h.adminDeploymentStream(ctx, filter, snapshot, subscription)
	cache, buffering := "no-cache, no-store", "no"
	return api.StreamAdminDeployments200TexteventStreamResponse{
		Body: body, Headers: api.StreamAdminDeployments200ResponseHeaders{CacheControl: &cache, XAccelBuffering: &buffering},
	}, nil
}

func (h *Handler) contestantDeploymentStream(ctx context.Context, team int64, problem string, snapshot []api.ContestantDeployment, subscription service.Subscription) io.Reader {
	reader, writer := io.Pipe()
	go func() {
		defer closeDeploymentSubscription(subscription)
		defer writer.Close()
		if err := writeSSE(writer, "snapshot", uuid.NewString(), api.ContestantDeploymentsResponse{Deployments: snapshot}); err != nil {
			return
		}
		h.streamDeploymentEvents(ctx, writer, subscription, func(deployment core.Deployment) (string, any, bool) {
			if deployment.TeamCode != team || deployment.ProblemCode != problem {
				return "", nil, false
			}
			mapped, err := h.toContestantDeployment(ctx, deployment)
			if err != nil {
				return "", nil, false
			}
			return latestDeploymentEventID(deployment), api.ContestantDeploymentResponse{Deployment: mapped}, true
		})
	}()
	return reader
}

func (h *Handler) adminDeploymentStream(ctx context.Context, filter core.DeploymentFilter, snapshot []api.AdminDeployment, subscription service.Subscription) io.Reader {
	reader, writer := io.Pipe()
	go func() {
		defer closeDeploymentSubscription(subscription)
		defer writer.Close()
		if err := writeSSE(writer, "snapshot", uuid.NewString(), api.AdminDeploymentsResponse{Deployments: snapshot}); err != nil {
			return
		}
		h.streamDeploymentEvents(ctx, writer, subscription, func(deployment core.Deployment) (string, any, bool) {
			if filter.TeamCode != nil && deployment.TeamCode != *filter.TeamCode ||
				filter.ProblemCode != nil && deployment.ProblemCode != *filter.ProblemCode {
				return "", nil, false
			}
			return latestDeploymentEventID(deployment), api.AdminDeploymentResponse{Deployment: toAPIAdminDeployment(deployment)}, true
		})
	}()
	return reader
}

func (h *Handler) streamDeploymentEvents(ctx context.Context, writer io.Writer, subscription service.Subscription, mapEvent func(core.Deployment) (string, any, bool)) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	if subscription == nil {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
					return
				}
			}
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
				return
			}
		case payload, ok := <-subscription.Messages():
			if !ok {
				return
			}
			var deployment core.Deployment
			if json.Unmarshal(payload, &deployment) != nil {
				continue
			}
			team, problem := deployment.TeamCode, deployment.ProblemCode
			latest, err := h.service.Store.ListDeployments(ctx, core.DeploymentFilter{TeamCode: &team, ProblemCode: &problem})
			if err != nil {
				continue
			}
			found := false
			for _, current := range latest {
				if current.RequestID == deployment.RequestID {
					deployment, found = current, true
					break
				}
			}
			if !found {
				continue
			}
			id, body, include := mapEvent(deployment)
			if include && writeSSE(writer, "deployment", id, body) != nil {
				return
			}
		}
	}
}

func writeSSE(writer io.Writer, event, id string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "event: %s\nid: %s\ndata: %s\n\n", event, id, encoded)
	return err
}

func (h *Handler) toContestantDeployment(ctx context.Context, deployment core.Deployment) (api.ContestantDeployment, error) {
	snapshot, err := h.service.Store.GetContent(ctx, deployment.ContentCommit)
	if err != nil {
		return api.ContestantDeployment{}, err
	}
	problem, ok := snapshot.Problem(deployment.ProblemCode)
	if !ok {
		return api.ContestantDeployment{}, core.NewError(http.StatusInternalServerError, "internal_error", "Deployment references missing problem snapshot")
	}
	penalty := core.DeploymentPenalty(problem.MaxScore, problem.Redeploy, deployment.Revision)
	allowed := int32(-1)
	if problem.Redeploy.Type == core.RedeployUnredeployable {
		allowed = 0
	}
	return toAPIContestantDeployment(deployment, penalty, allowed), nil
}
