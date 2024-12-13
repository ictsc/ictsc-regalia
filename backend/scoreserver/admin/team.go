package admin

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"github.com/jmoiron/sqlx"
)

type TeamServiceHandler struct {
	ListWorkflow   domain.TeamListWorkflow
	GetWorkflow    domain.TeamGetWorkflow
	CreateWorkflow domain.TeamCreateWorkflow
	UpdateWorkflow domain.TeamUpdateWorkflow
	DeleteWorkflow domain.TeamDeleteWorkflow

	adminv1connect.UnimplementedTeamServiceHandler
}

func NewTeamServiceHandler(db *sqlx.DB) *TeamServiceHandler {
	repo := pg.NewRepository(db)

	return &TeamServiceHandler{
		ListWorkflow: domain.TeamListWorkflow{Lister: repo},
		GetWorkflow:  domain.TeamGetWorkflow{Getter: repo},
		CreateWorkflow: domain.TeamCreateWorkflow{
			RunTx: func(ctx context.Context, f func(domain.TeamCreateTxEffect) error) error {
				return repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
					return f(tx)
				})
			},
		},
		UpdateWorkflow: domain.TeamUpdateWorkflow{
			RunTx: func(ctx context.Context, f func(domain.TeamUpdateTxEffect) error) error {
				return repo.RunTx(ctx, func(tx *pg.RepositoryTx) error { return f(tx) })
			},
		},
		DeleteWorkflow: domain.TeamDeleteWorkflow{
			RunTx: func(ctx context.Context, f func(domain.TeamDeleteTxEffect) error) error {
				return repo.RunTx(ctx, func(tx *pg.RepositoryTx) error {
					return f(tx)
				})
			},
		},

		UnimplementedTeamServiceHandler: adminv1connect.UnimplementedTeamServiceHandler{},
	}
}

func (h *TeamServiceHandler) ListTeams(
	ctx context.Context,
	req *connect.Request[adminv1.ListTeamsRequest],
) (*connect.Response[adminv1.ListTeamsResponse], error) {
	teams, err := h.ListWorkflow.Run(ctx)
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
	team, err := h.GetWorkflow.Run(ctx, domain.TeamGetInput{
		Code: int(req.Msg.GetCode()),
	})
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
	team, err := h.CreateWorkflow.Run(ctx, domain.TeamCreateInput{
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
	protoTeam := req.Msg.GetTeam()

	team, err := h.UpdateWorkflow.Run(ctx, domain.TeamUpdateInput{
		Code:         int(protoTeam.GetCode()),
		Name:         protoTeam.GetName(),
		Organization: protoTeam.GetOrganization(),
	})
	if err != nil {
		return nil, connectError(err)
	}

	return connect.NewResponse(&adminv1.UpdateTeamResponse{
		Team: convertTeam(team),
	}), nil
}

func (h *TeamServiceHandler) DeleteTeam(
	ctx context.Context,
	req *connect.Request[adminv1.DeleteTeamRequest],
) (*connect.Response[adminv1.DeleteTeamResponse], error) {
	if err := h.DeleteWorkflow.Run(ctx, domain.TeamDeleteInput{
		Code: int(req.Msg.GetCode()),
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

func connectError(err error) error {
	if err == nil {
		return nil
	}
	var code connect.Code
	switch domain.ErrTypeFrom(err) {
	case domain.ErrTypeInternal:
		code = connect.CodeInternal
	case domain.ErrTypeInvalidArgument:
		code = connect.CodeInvalidArgument
	case domain.ErrTypeAlreadyExists:
		code = connect.CodeAlreadyExists
	case domain.ErrTypeNotFound:
		code = connect.CodeNotFound
	case domain.ErrTypeUnknown:
		fallthrough
	default:
		code = connect.CodeUnknown
	}
	return connect.NewError(code, err)
}
