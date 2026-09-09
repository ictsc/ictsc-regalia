package httpserver

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func (h *Handler) ListAdminTeams(ctx context.Context, _ api.ListAdminTeamsRequestObject) (api.ListAdminTeamsResponseObject, error) {
	teams, err := h.service.Store.ListTeams(ctx)
	if err != nil {
		return nil, err
	}
	mapped := make([]api.Team, 0, len(teams))
	for _, team := range teams {
		mapped = append(mapped, toAPITeam(team))
	}
	return api.ListAdminTeams200JSONResponse{Teams: mapped}, nil
}

func (h *Handler) CreateAdminTeam(ctx context.Context, request api.CreateAdminTeamRequestObject) (api.CreateAdminTeamResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	color := core.DefaultTeamColor
	if request.Body.Color != nil {
		color = string(*request.Body.Color)
	}
	team, err := h.service.Store.CreateTeam(ctx, core.Team{
		Color: color,
		Code:  request.Body.Code, Name: request.Body.Name, Organization: request.Body.Organization, MemberLimit: request.Body.MemberLimit,
	})
	if err != nil {
		return nil, err
	}
	return api.CreateAdminTeam201JSONResponse{Team: toAPITeam(team)}, nil
}

func (h *Handler) GetAdminTeam(ctx context.Context, request api.GetAdminTeamRequestObject) (api.GetAdminTeamResponseObject, error) {
	team, err := h.service.Store.GetTeam(ctx, request.TeamCode)
	if err != nil {
		return nil, err
	}
	return api.GetAdminTeam200JSONResponse{Team: toAPITeam(team)}, nil
}

func (h *Handler) UpdateAdminTeam(ctx context.Context, request api.UpdateAdminTeamRequestObject) (api.UpdateAdminTeamResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	var color *string
	if request.Body.Color != nil {
		value := string(*request.Body.Color)
		color = &value
	}
	team, err := h.service.Store.UpdateTeam(ctx, request.TeamCode, core.TeamPatch{
		Color: color,
		Name:  request.Body.Name, Organization: request.Body.Organization, MemberLimit: request.Body.MemberLimit,
	})
	if err != nil {
		return nil, err
	}
	return api.UpdateAdminTeam200JSONResponse{Team: toAPITeam(team)}, nil
}

func (h *Handler) DeleteAdminTeam(ctx context.Context, request api.DeleteAdminTeamRequestObject) (api.DeleteAdminTeamResponseObject, error) {
	if err := h.service.Store.DeleteTeam(ctx, request.TeamCode); err != nil {
		return nil, err
	}
	return api.DeleteAdminTeam204Response{}, nil
}

func (h *Handler) ListAdminInvitations(ctx context.Context, request api.ListAdminInvitationsRequestObject) (api.ListAdminInvitationsResponseObject, error) {
	invitations, err := h.service.Store.ListInvitations(ctx)
	if err != nil {
		return nil, err
	}
	now := h.service.Now()
	mapped := make([]api.Invitation, 0)
	for _, invitation := range invitations {
		if invitation.TeamCode != int64(request.Params.TeamCode) {
			continue
		}
		if (request.Params.IncludeExpired == nil || !*request.Params.IncludeExpired) && !now.Before(invitation.ExpiresAt) {
			continue
		}
		mapped = append(mapped, toAPIInvitation(invitation))
	}
	return api.ListAdminInvitations200JSONResponse{Invitations: mapped}, nil
}

func (h *Handler) CreateAdminInvitation(ctx context.Context, request api.CreateAdminInvitationRequestObject) (api.CreateAdminInvitationResponseObject, error) {
	if request.Body == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "Request body is required")
	}
	now := h.service.Now()
	if !request.Body.ExpiresAt.After(now) {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "expires_at must be in the future")
	}
	if _, err := h.service.Store.GetTeam(ctx, request.Body.TeamCode); err != nil {
		return nil, err
	}
	invitation, err := h.service.Store.CreateInvitation(ctx, core.Invitation{
		Code: uuid.NewString(), TeamCode: request.Body.TeamCode, CreatedAt: now, ExpiresAt: request.Body.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	return api.CreateAdminInvitation201JSONResponse{Invitation: toAPIInvitation(invitation)}, nil
}

func (h *Handler) ListAdminContestants(ctx context.Context, request api.ListAdminContestantsRequestObject) (api.ListAdminContestantsResponseObject, error) {
	contestants, err := h.service.Store.ListContestants(ctx)
	if err != nil {
		return nil, err
	}
	mapped := make([]api.AdminContestant, 0, len(contestants))
	for _, contestant := range contestants {
		if request.Params.TeamCode != nil && contestant.TeamCode != int64(*request.Params.TeamCode) {
			continue
		}
		team, err := h.service.Store.GetTeam(ctx, contestant.TeamCode)
		if err != nil {
			return nil, err
		}
		mapped = append(mapped, api.AdminContestant{Profile: toAPIProfile(contestant), Team: toAPITeam(team), DiscordId: contestant.DiscordID})
	}
	return api.ListAdminContestants200JSONResponse{Contestants: mapped}, nil
}

func (h *Handler) GetContestantProfile(ctx context.Context, _ api.GetContestantProfileRequestObject) (api.GetContestantProfileResponseObject, error) {
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetContestantProfile200JSONResponse{Profile: toAPIProfile(contestant)}, nil
}

func (h *Handler) UpdateContestantProfile(ctx context.Context, request api.UpdateContestantProfileRequestObject) (api.UpdateContestantProfileResponseObject, error) {
	if request.Body == nil || request.Body.DisplayName == nil && request.Body.SelfIntroduction == nil {
		return nil, core.NewError(http.StatusUnprocessableEntity, "validation_error", "At least one profile field is required")
	}
	contestant, err := h.currentContestant(ctx)
	if err != nil {
		return nil, err
	}
	contestant, err = h.service.Store.UpdateContestant(ctx, contestant.Name, core.ContestantPatch{
		DisplayName: request.Body.DisplayName, SelfIntroduction: request.Body.SelfIntroduction,
	})
	if err != nil {
		return nil, err
	}
	return api.UpdateContestantProfile200JSONResponse{Profile: toAPIProfile(contestant)}, nil
}

func (h *Handler) ListContestantTeams(ctx context.Context, _ api.ListContestantTeamsRequestObject) (api.ListContestantTeamsResponseObject, error) {
	teams, err := h.service.Store.ListTeams(ctx)
	if err != nil {
		return nil, err
	}
	contestants, err := h.service.Store.ListContestants(ctx)
	if err != nil {
		return nil, err
	}
	mapped := make([]api.TeamProfile, 0, len(teams))
	for _, team := range teams {
		members := make([]api.ContestantProfile, 0)
		for _, contestant := range contestants {
			if contestant.TeamCode == team.Code {
				members = append(members, toAPIProfile(contestant))
			}
		}
		sort.Slice(members, func(i, j int) bool { return members[i].Name < members[j].Name })
		mapped = append(mapped, api.TeamProfile{Team: toAPITeam(team), Members: members})
	}
	return api.ListContestantTeams200JSONResponse{Teams: mapped}, nil
}

func (h *Handler) currentContestant(ctx context.Context) (core.Contestant, error) {
	name := principal(ctx).Session.ContestantName
	if name == "" {
		return core.Contestant{}, core.NewError(http.StatusUnauthorized, "authentication_required", "Contestant authentication is required")
	}
	return h.service.Store.GetContestant(ctx, name)
}
