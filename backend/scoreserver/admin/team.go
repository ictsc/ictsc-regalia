package admin

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type TeamServiceHandler struct {
	Enforcer     *auth.Enforcer
	ListEffect   domain.TeamListEffect
	GetEffect    domain.TeamGetEffect
	CreateEffect domain.TeamCreateEffect
	UpdateEffect domain.TeamUpdateEffect
	DeleteEffect domain.Tx[teamDeleteEffect]
}

var _ adminv1connect.TeamServiceHandler = (*TeamServiceHandler)(nil)

func NewTeamServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *TeamServiceHandler {
	return &TeamServiceHandler{
		Enforcer: enforcer,

		ListEffect:   repo,
		GetEffect:    repo,
		CreateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.TeamCreateTxEffect { return rt }),
		UpdateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.TeamUpdateTxEffect { return rt }),
		DeleteEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) teamDeleteEffect { return rt }),
	}
}

func (h *TeamServiceHandler) ListTeams(
	ctx context.Context,
	req *connect.Request[adminv1.ListTeamsRequest],
) (*connect.Response[adminv1.ListTeamsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "teams", "list"); err != nil {
		return nil, err
	}

	teams, err := domain.ListTeams(ctx, h.ListEffect)
	if err != nil {
		return nil, connectError(err)
	}

	protoTeams := make([]*adminv1.Team, 0, len(teams))
	for _, team := range teams {
		protoTeams = append(protoTeams, convertTeam(team))
	}

	return connect.NewResponse(&adminv1.ListTeamsResponse{
		Teams: protoTeams,
	}), nil
}

func (h *TeamServiceHandler) GetTeam(
	ctx context.Context,
	req *connect.Request[adminv1.GetTeamRequest],
) (*connect.Response[adminv1.GetTeamResponse], error) {
	if err := enforce(ctx, h.Enforcer, "teams", "get"); err != nil {
		return nil, err
	}

	inCode := req.Msg.GetCode()
	if inCode == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("code is required"))
	}

	code, err := domain.NewTeamCode(inCode)
	if err != nil {
		return nil, connectError(err)
	}

	team, err := code.Team(ctx, h.GetEffect)
	if err != nil {
		return nil, connectError(err)
	}

	return connect.NewResponse(&adminv1.GetTeamResponse{
		Team: convertTeam(team),
	}), nil
}

func (h *TeamServiceHandler) CreateTeam(
	ctx context.Context,
	req *connect.Request[adminv1.CreateTeamRequest],
) (*connect.Response[adminv1.CreateTeamResponse], error) {
	if err := enforce(ctx, h.Enforcer, "teams", "create"); err != nil {
		return nil, err
	}
	team, err := domain.CreateTeam(ctx, h.CreateEffect, domain.TeamCreateInput{
		Code:         int(req.Msg.GetTeam().GetCode()),
		Name:         req.Msg.GetTeam().GetName(),
		Organization: req.Msg.GetTeam().GetOrganization(),
	})
	if err != nil {
		return nil, connectError(err)
	}

	return connect.NewResponse(&adminv1.CreateTeamResponse{
		Team: convertTeam(team),
	}), nil
}

func (h *TeamServiceHandler) UpdateTeam(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateTeamRequest],
) (*connect.Response[adminv1.UpdateTeamResponse], error) {
	if err := enforce(ctx, h.Enforcer, "teams", "update"); err != nil {
		return nil, err
	}

	protoTeam := req.Msg.GetTeam()

	protoCode := protoTeam.GetCode()
	if protoCode == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("code is required"))
	}

	teamCode, err := domain.NewTeamCode(protoCode)
	if err != nil {
		return nil, connectError(err)
	}

	team, err := teamCode.Team(ctx, h.GetEffect)
	if err != nil {
		return nil, connectError(err)
	}

	if name := protoTeam.GetName(); name != "" && name != team.Name() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be updated"))
	}

	if err := team.Update(ctx, h.UpdateEffect, domain.TeamUpdateInput{
		Organization: protoTeam.GetOrganization(),
	}); err != nil {
		return nil, connectError(err)
	}

	return connect.NewResponse(&adminv1.UpdateTeamResponse{
		Team: convertTeam(team),
	}), nil
}

type teamDeleteEffect interface {
	domain.TeamGetter
	domain.TeamDeleter
}

func (h *TeamServiceHandler) DeleteTeam(
	ctx context.Context,
	req *connect.Request[adminv1.DeleteTeamRequest],
) (*connect.Response[adminv1.DeleteTeamResponse], error) {
	if err := enforce(ctx, h.Enforcer, "teams", "delete"); err != nil {
		return nil, err
	}

	protoCode := req.Msg.GetCode()
	if protoCode == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("code is required"))
	}
	code, err := domain.NewTeamCode(protoCode)
	if err != nil {
		return nil, connectError(err)
	}

	if err := h.DeleteEffect.RunInTx(ctx, func(effect teamDeleteEffect) error {
		team, err := code.Team(ctx, effect)
		if err != nil {
			return err
		}

		return team.Delete(ctx, effect)
	}); err != nil {
		return nil, connectError(err)
	}

	return connect.NewResponse(&adminv1.DeleteTeamResponse{}), nil
}

func convertTeam(team *domain.Team) *adminv1.Team {
	return &adminv1.Team{
		Code:         int64(team.Code()),
		Name:         team.Name(),
		Organization: team.Organization(),
	}
}
