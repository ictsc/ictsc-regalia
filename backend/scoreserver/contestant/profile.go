package contestant

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type ProfileServiceHandler struct {
	contestantv1connect.UnimplementedProfileServiceHandler

	ListTeamsEffect ListTeamsEffect
}

var _ contestantv1connect.ProfileServiceHandler = (*ProfileServiceHandler)(nil)

func newProfileServiceHandler(repo *pg.Repository) *ProfileServiceHandler {
	return &ProfileServiceHandler{
		ListTeamsEffect: repo,
	}
}

type ListTeamsEffect interface {
	domain.TeamsLister
	domain.TeamMemberProfileReader
}

func (h *ProfileServiceHandler) ListTeams(
	ctx context.Context,
	req *connect.Request[contestantv1.ListTeamsRequest],
) (*connect.Response[contestantv1.ListTeamsResponse], error) {
	if _, err := session.UserSessionStore.Get(ctx); err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	teams, err := domain.ListTeams(ctx, h.ListTeamsEffect)
	if err != nil {
		return nil, err
	}
	profileByTeam := make(map[domain.TeamCode][]*contestantv1.ContestantProfile)
	for _, team := range teams {
		profileByTeam[team.Code()] = []*contestantv1.ContestantProfile{}
	}
	members, err := domain.ListTeamMembers(ctx, h.ListTeamsEffect)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		code := member.Team().Code()
		profileByTeam[code] = append(profileByTeam[code], &contestantv1.ContestantProfile{
			Name:        string(member.Name()),
			DisplayName: member.UserProfile().DisplayName(),
		})
	}

	teamProfiles := make([]*contestantv1.TeamProfile, 0, len(teams))
	for _, team := range teams {
		teamProfiles = append(teamProfiles, &contestantv1.TeamProfile{
			Name:         team.Name(),
			Organization: team.Organization(),
			Members:      profileByTeam[team.Code()],
		})
	}

	return connect.NewResponse(&contestantv1.ListTeamsResponse{
		Teams: teamProfiles,
	}), nil
}
