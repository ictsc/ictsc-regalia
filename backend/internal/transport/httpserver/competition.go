package httpserver

import (
	"context"
	"errors"
	"net/http"
	"sort"

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
	ranking, err := h.service.Ranking(ctx, true)
	if err != nil {
		return nil, err
	}
	return api.GetAdminRanking200JSONResponse(toAPIRanking(ranking)), nil
}

func (h *Handler) GetContestantRanking(ctx context.Context, _ api.GetContestantRankingRequestObject) (api.GetContestantRankingResponseObject, error) {
	ranking, err := h.service.Ranking(ctx, false)
	if err != nil {
		return nil, err
	}
	return api.GetContestantRanking200JSONResponse(toAPIRanking(ranking)), nil
}

func (h *Handler) GetAdminRule(ctx context.Context, _ api.GetAdminRuleRequestObject) (api.GetAdminRuleResponseObject, error) {
	markdown, err := h.ruleMarkdown(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetAdminRule200JSONResponse{Rule: api.Rule{Markdown: markdown}}, nil
}

func (h *Handler) GetContestantRule(ctx context.Context, _ api.GetContestantRuleRequestObject) (api.GetContestantRuleResponseObject, error) {
	markdown, err := h.ruleMarkdown(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetContestantRule200JSONResponse{Rule: api.Rule{Markdown: markdown}}, nil
}

// ruleMarkdown is the active content's rule.md; the rule is managed in the
// content repository only. Before the first import it is empty.
func (h *Handler) ruleMarkdown(ctx context.Context) (string, error) {
	snapshot, err := h.service.Store.ActiveContent(ctx)
	var domainErr *core.Error
	if errors.As(err, &domainErr) && domainErr.Code == "content_not_available" {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return snapshot.Manifest.RuleMarkdown, nil
}

func (h *Handler) GetAdminDashboardSchedule(ctx context.Context, _ api.GetAdminDashboardScheduleRequestObject) (api.GetAdminDashboardScheduleResponseObject, error) {
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetAdminDashboardSchedule200JSONResponse(h.toAPIDashboardSchedule(state)), nil
}

func (h *Handler) GetContestantDashboardSchedule(ctx context.Context, _ api.GetContestantDashboardScheduleRequestObject) (api.GetContestantDashboardScheduleResponseObject, error) {
	state, err := h.service.Store.GetCompetitionState(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetContestantDashboardSchedule200JSONResponse(h.toAPIDashboardSchedule(state)), nil
}

func (h *Handler) ReplaceAdminDashboardSchedule(ctx context.Context, request api.ReplaceAdminDashboardScheduleRequestObject) (api.ReplaceAdminDashboardScheduleResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	if request.Body.RankingFreezeAt != nil && h.service.Config.SecondaryQualifierMode {
		return nil, core.NewError(http.StatusConflict, "conflict", "Ranking freeze is not used in secondary qualifier mode")
	}
	state, err := h.service.Store.ReplaceFreezeAt(ctx, request.Body.RankingFreezeAt, principal(ctx).Session.AdminName)
	if err != nil {
		return nil, err
	}
	return api.ReplaceAdminDashboardSchedule200JSONResponse(h.toAPIDashboardSchedule(state)), nil
}

func (h *Handler) toAPIDashboardSchedule(state core.CompetitionState) api.DashboardScheduleResponse {
	return api.DashboardScheduleResponse{
		DashboardSchedule: api.DashboardSchedule{RankingFreezeAt: state.RankingFreezeAt},
		CompetitionMode:   api.CompetitionMode(h.service.CompetitionMode()),
	}
}

func (h *Handler) GetAdminDevContestPhase(ctx context.Context, _ api.GetAdminDevContestPhaseRequestObject) (api.GetAdminDevContestPhaseResponseObject, error) {
	phase, sections, err := h.service.ContestPhase(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetAdminDevContestPhase200JSONResponse{ContestPhase: h.toAPIDevContestPhase(phase, sections)}, nil
}

func (h *Handler) ReplaceAdminDevContestPhase(ctx context.Context, request api.ReplaceAdminDevContestPhaseRequestObject) (api.ReplaceAdminDevContestPhaseResponseObject, error) {
	if request.Body == nil || !request.Body.Phase.Valid() {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "phase must be BEFORE, DURING or AFTER")
	}
	phase, sections, err := h.service.SetDevContestPhase(ctx, core.ContestPhase(request.Body.Phase))
	if err != nil {
		return nil, err
	}
	return api.ReplaceAdminDevContestPhase200JSONResponse{ContestPhase: h.toAPIDevContestPhase(phase, sections)}, nil
}

func (h *Handler) toAPIDevContestPhase(phase core.ContestPhase, sections []core.Section) api.DevContestPhase {
	mapped := make([]api.ContestantScheduleSection, 0, len(sections))
	for _, section := range sections {
		mapped = append(mapped, api.ContestantScheduleSection{Slug: section.Slug, Beginning: section.Beginning, Ending: section.Ending})
	}
	sort.SliceStable(mapped, func(i, j int) bool {
		if mapped[i].Beginning.Equal(mapped[j].Beginning) {
			return mapped[i].Slug < mapped[j].Slug
		}
		return mapped[i].Beginning.Before(mapped[j].Beginning)
	})
	return api.DevContestPhase{Enabled: h.service.DevContestPhaseEnabled(), Phase: api.ContestPhase(phase), Sections: mapped}
}
