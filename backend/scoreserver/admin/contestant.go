package admin

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type ContestantServiceHandler struct {
	adminv1connect.UnimplementedContestantServiceHandler

	Enforcer   *auth.Enforcer
	ListEffect ListContestantsEffect
}

var _ adminv1connect.ContestantServiceHandler = (*ContestantServiceHandler)(nil)

func newContestantServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *ContestantServiceHandler {
	return &ContestantServiceHandler{
		Enforcer:   enforcer,
		ListEffect: repo,
	}
}

type ListContestantsEffect interface {
	domain.TeamGetter
	domain.TeamMemberProfileReader
}

func (h *ContestantServiceHandler) ListContestants(
	ctx context.Context,
	req *connect.Request[adminv1.ListContestantsRequest],
) (*connect.Response[adminv1.ListContestantsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "contestants", "list"); err != nil {
		return nil, err
	}

	protoCode := req.Msg.GetTeamCode()

	members, err := listMembers(ctx, protoCode, h.ListEffect)
	if err != nil {
		return nil, err
	}

	contestants := make([]*adminv1.Contestant, 0, len(members))
	for _, member := range members {
		contestants = append(contestants, &adminv1.Contestant{
			Name:        string(member.UserProfile().Name()),
			DisplayName: member.UserProfile().DisplayName(),
			Team:        convertTeam(member.Team()),
			Profile:     &adminv1.Profile{},
			DiscordId:   member.DiscordUserID().String(),
		})
	}

	return connect.NewResponse(&adminv1.ListContestantsResponse{
		Contestants: contestants,
	}), nil
}

func listMembers(ctx context.Context, code int64, eff ListContestantsEffect) ([]*domain.TeamMemberProfile, error) {
	if code == 0 {
		return domain.ListTeamMembers(ctx, eff)
	}

	teamCode, err := domain.NewTeamCode(code)
	if err != nil {
		return nil, err
	}

	team, err := teamCode.Team(ctx, eff)
	if err != nil {
		return nil, err
	}

	return team.Members(ctx, eff)
}
