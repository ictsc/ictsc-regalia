package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
	v1 "github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1"
	"github.com/ictsc/ictsc-outlands/backend/internal/proto/anita/v1/v1connect"
	"github.com/ictsc/ictsc-outlands/backend/pkg/optional"
)

// TeamServiceHandler チームサービスのハンドラ
type TeamServiceHandler struct {
	srv service.TeamService
}

var _ v1connect.TeamServiceHandler = (*TeamServiceHandler)(nil)

// NewTeamServiceHandler チームサービスのハンドラを作成する
func NewTeamServiceHandler(srv service.TeamService) *TeamServiceHandler {
	return &TeamServiceHandler{srv: srv}
}

// GetTeam チームを取得する
func (s *TeamServiceHandler) GetTeam(ctx context.Context, req *connect.Request[v1.GetTeamRequest]) (*connect.Response[v1.GetTeamResponse], error) {
	id, err := value.NewTeamID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	team, err := s.srv.ReadTeam(ctx, id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetTeamResponse{
		Team: fromDomainTeam(team),
	}), nil
}

// GetTeams チーム一覧を取得する
func (s *TeamServiceHandler) GetTeams(ctx context.Context, _ *connect.Request[v1.GetTeamsRequest]) (*connect.Response[v1.GetTeamsResponse], error) {
	teams, err := s.srv.ReadTeams(ctx)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetTeamsResponse{
		Teams: fromDomainTeams(teams),
	}), nil
}

// GetConnectionInfo チームの接続情報を取得する
func (s *TeamServiceHandler) GetConnectionInfo(
	ctx context.Context, req *connect.Request[v1.GetConnectionInfoRequest],
) (*connect.Response[v1.GetConnectionInfoResponse], error) {
	id, err := value.NewTeamID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	team, err := s.srv.ReadTeam(ctx, id)
	if err != nil {
		return nil, err
	}

	var bastion *v1.Bastion
	if team.Bastion().Valid {
		bastion = fromDomainBastion(team.Bastion().V)
	}

	return connect.NewResponse(&v1.GetConnectionInfoResponse{
		Bastion: bastion,
	}), nil
}

// GetMembers チームのメンバー一覧を取得する
func (s *TeamServiceHandler) GetMembers(
	ctx context.Context, req *connect.Request[v1.GetMembersRequest],
) (*connect.Response[v1.GetMembersResponse], error) {
	id, err := value.NewTeamID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	team, err := s.srv.ReadTeam(ctx, id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetMembersResponse{
		Members: fromDomainUsers(team.Members()),
	}), nil
}

// PatchTeam チームを更新する
func (s *TeamServiceHandler) PatchTeam(
	ctx context.Context, req *connect.Request[v1.PatchTeamRequest],
) (*connect.Response[v1.PatchTeamResponse], error) {
	id, err := value.NewTeamID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	args := service.UpdateTeamArgs{
		ID:            id,
		Code:          optional.New(value.TeamCode{}, false),
		Name:          optional.New(value.TeamName{}, false),
		Org:           optional.New(value.TeamOrganization{}, false),
		CodeRemaining: optional.New(value.TeamCodeRemaining{}, false),
		Bastion:       optional.New(value.Bastion{}, false),
	}

	if req.Msg.Code != nil {
		code, err := value.NewTeamCode(int(req.Msg.GetCode()))
		if err != nil {
			return nil, err
		}

		args.Code = optional.NewValid(code)
	}

	if req.Msg.Name != nil {
		name, err := value.NewTeamName(req.Msg.GetName())
		if err != nil {
			return nil, err
		}

		args.Name = optional.NewValid(name)
	}

	if req.Msg.Organization != nil {
		org, err := value.NewTeamOrganization(req.Msg.GetOrganization())
		if err != nil {
			return nil, err
		}

		args.Org = optional.NewValid(org)
	}

	if req.Msg.CodeRemaining != nil {
		codeRemaining, err := value.NewTeamCodeRemaining(int(req.Msg.GetCodeRemaining()))
		if err != nil {
			return nil, err
		}

		args.CodeRemaining = optional.NewValid(codeRemaining)
	}

	team, err := s.srv.UpdateTeam(ctx, args)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.PatchTeamResponse{
		Team: fromDomainTeam(team),
	}), nil
}

// PutConnectionInfo チームの接続情報を更新する
func (s *TeamServiceHandler) PutConnectionInfo(
	ctx context.Context, req *connect.Request[v1.PutConnectionInfoRequest],
) (*connect.Response[v1.PutConnectionInfoResponse], error) {
	id, err := value.NewTeamID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	bastion, err := value.NewBastion(
		req.Msg.GetBastion().GetUser(),
		req.Msg.GetBastion().GetPassword(),
		req.Msg.GetBastion().GetHost(),
		int(req.Msg.GetBastion().GetPort()),
	)
	if err != nil {
		return nil, err
	}

	args := service.UpdateTeamArgs{
		ID:            id,
		Code:          optional.New(value.TeamCode{}, false),
		Name:          optional.New(value.TeamName{}, false),
		Org:           optional.New(value.TeamOrganization{}, false),
		CodeRemaining: optional.New(value.TeamCodeRemaining{}, false),
		Bastion:       optional.NewValid(bastion),
	}

	_, err = s.srv.UpdateTeam(ctx, args)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.PutConnectionInfoResponse{}), nil
}

// PostTeam チームを作成する
func (s *TeamServiceHandler) PostTeam(ctx context.Context, req *connect.Request[v1.PostTeamRequest]) (*connect.Response[v1.PostTeamResponse], error) {
	code, err := value.NewTeamCode(int(req.Msg.GetCode()))
	if err != nil {
		return nil, err
	}

	name, err := value.NewTeamName(req.Msg.GetName())
	if err != nil {
		return nil, err
	}

	org, err := value.NewTeamOrganization(req.Msg.GetOrganization())
	if err != nil {
		return nil, err
	}

	codeRemaining, err := value.NewTeamCodeRemaining(int(req.Msg.GetCodeRemaining()))
	if err != nil {
		return nil, err
	}

	args := service.CreateTeamArgs{
		Code:          code,
		Name:          name,
		Org:           org,
		CodeRemaining: codeRemaining,
	}

	team, err := s.srv.CreateTeam(ctx, args)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.PostTeamResponse{
		Team: fromDomainTeam(team),
	}), nil
}

// DeleteTeam チームを削除する
func (s *TeamServiceHandler) DeleteTeam(
	ctx context.Context, req *connect.Request[v1.DeleteTeamRequest],
) (*connect.Response[v1.DeleteTeamResponse], error) {
	id, err := value.NewTeamID(req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	err = s.srv.DeleteTeam(ctx, id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteTeamResponse{}), nil
}

// MoveMember チームのメンバーを移動する
func (s *TeamServiceHandler) MoveMember(
	ctx context.Context, req *connect.Request[v1.MoveMemberRequest],
) (*connect.Response[v1.MoveMemberResponse], error) {
	to, err := value.NewTeamID(req.Msg.GetToTeamId())
	if err != nil {
		return nil, err
	}

	memberID, err := value.NewUserID(req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}

	err = s.srv.MoveMember(ctx, to, memberID)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.MoveMemberResponse{}), nil
}
