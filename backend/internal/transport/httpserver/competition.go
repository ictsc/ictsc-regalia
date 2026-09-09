package httpserver

import (
	"context"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) GetHealth(ctx context.Context, _ api.GetHealthRequestObject) (api.GetHealthResponseObject, error) {
	if err := h.service.Store.Ping(ctx); err != nil {
		return nil, err
	}
	return api.GetHealth200JSONResponse{Status: api.Ok}, nil
}

func (h *Handler) GetAdminRanking(ctx context.Context, _ api.GetAdminRankingRequestObject) (api.GetAdminRankingResponseObject, error) {
	ranking, frozen, frozenAt, err := h.service.Ranking(ctx, true)
	if err != nil {
		return nil, err
	}
	return api.GetAdminRanking200JSONResponse(toAPIRanking(ranking, frozen, frozenAt)), nil
}

func (h *Handler) GetContestantRanking(ctx context.Context, _ api.GetContestantRankingRequestObject) (api.GetContestantRankingResponseObject, error) {
	ranking, frozen, frozenAt, err := h.service.Ranking(ctx, false)
	if err != nil {
		return nil, err
	}
	return api.GetContestantRanking200JSONResponse(toAPIRanking(ranking, frozen, frozenAt)), nil
}

func (h *Handler) GetAdminRule(ctx context.Context, _ api.GetAdminRuleRequestObject) (api.GetAdminRuleResponseObject, error) {
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetAdminRule200JSONResponse{Rule: api.Rule{Markdown: state.RuleMarkdown}}, nil
}

func (h *Handler) GetContestantRule(ctx context.Context, _ api.GetContestantRuleRequestObject) (api.GetContestantRuleResponseObject, error) {
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetContestantRule200JSONResponse{Rule: api.Rule{Markdown: state.RuleMarkdown}}, nil
}

func (h *Handler) ReplaceAdminRule(ctx context.Context, request api.ReplaceAdminRuleRequestObject) (api.ReplaceAdminRuleResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	state, err := h.service.Store.ReplaceRule(ctx, request.Body.Markdown, principal(ctx).Session.AdminName)
	if err != nil {
		return nil, err
	}
	return api.ReplaceAdminRule200JSONResponse{Rule: api.Rule{Markdown: state.RuleMarkdown}}, nil
}

func (h *Handler) GetAdminDashboardSchedule(ctx context.Context, _ api.GetAdminDashboardScheduleRequestObject) (api.GetAdminDashboardScheduleResponseObject, error) {
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetAdminDashboardSchedule200JSONResponse{DashboardSchedule: api.DashboardSchedule{RankingFreezeAt: state.RankingFreezeAt}}, nil
}

func (h *Handler) GetContestantDashboardSchedule(ctx context.Context, _ api.GetContestantDashboardScheduleRequestObject) (api.GetContestantDashboardScheduleResponseObject, error) {
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetContestantDashboardSchedule200JSONResponse{DashboardSchedule: api.DashboardSchedule{RankingFreezeAt: state.RankingFreezeAt}}, nil
}

func (h *Handler) ReplaceAdminDashboardSchedule(ctx context.Context, request api.ReplaceAdminDashboardScheduleRequestObject) (api.ReplaceAdminDashboardScheduleResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	state, err := h.service.Store.ReplaceFreezeAt(ctx, request.Body.RankingFreezeAt, principal(ctx).Session.AdminName)
	if err != nil {
		return nil, err
	}
	return api.ReplaceAdminDashboardSchedule200JSONResponse{DashboardSchedule: api.DashboardSchedule{RankingFreezeAt: state.RankingFreezeAt}}, nil
}
