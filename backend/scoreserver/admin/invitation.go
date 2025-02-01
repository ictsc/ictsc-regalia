package admin

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvitationServiceHandler struct {
	adminv1connect.UnimplementedInvitationServiceHandler
	Enforcer       *auth.Enforcer
	ListWorkflow   *domain.InvitationCodeListWorkflow
	CreateWorkflow *domain.InvitationCodeCreateWorkflow
}

func NewInvitationServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *InvitationServiceHandler {
	return &InvitationServiceHandler{
		UnimplementedInvitationServiceHandler: adminv1connect.UnimplementedInvitationServiceHandler{},

		Enforcer:     enforcer,
		ListWorkflow: &domain.InvitationCodeListWorkflow{Lister: repo},
		CreateWorkflow: &domain.InvitationCodeCreateWorkflow{
			TeamGetter: repo,
			RunTx: func(ctx context.Context, f func(eff domain.InvitationCodeCreator) error) error {
				return repo.RunTx(ctx, func(tx *pg.RepositoryTx) error { return f(tx) })
			},
		},
	}
}

var _ adminv1connect.InvitationServiceHandler = (*InvitationServiceHandler)(nil)

func (h *InvitationServiceHandler) ListInvitationCodes(
	ctx context.Context,
	req *connect.Request[adminv1.ListInvitationCodesRequest],
) (*connect.Response[adminv1.ListInvitationCodesResponse], error) {
	if err := enforce(ctx, h.Enforcer, "invitation_codes", "list"); err != nil {
		return nil, err
	}

	ics, err := h.ListWorkflow.Run(ctx)
	if err != nil {
		return nil, connectError(err)
	}

	protoICs := make([]*adminv1.InvitationCode, 0, len(ics))
	for _, ic := range ics {
		protoICs = append(protoICs, convertInvitationCode(ic))
	}

	return connect.NewResponse(&adminv1.ListInvitationCodesResponse{
		InvitationCodes: protoICs,
	}), nil
}

func (h *InvitationServiceHandler) CreateInvitationCode(
	ctx context.Context,
	req *connect.Request[adminv1.CreateInvitationCodeRequest],
) (*connect.Response[adminv1.CreateInvitationCodeResponse], error) {
	if err := enforce(ctx, h.Enforcer, "invitation_codes", "create"); err != nil {
		return nil, err
	}

	code, err := h.CreateWorkflow.Run(ctx, domain.InvitationCodeCreateInput{
		TeamCode:  int(req.Msg.GetInvitationCode().GetTeamCode()),
		ExpiresAt: req.Msg.GetInvitationCode().GetExpiresAt().AsTime(),
	})
	if err != nil {
		return nil, connectError(err)
	}

	return connect.NewResponse(&adminv1.CreateInvitationCodeResponse{
		InvitationCode: convertInvitationCode(code),
	}), nil
}

func convertInvitationCode(ic *domain.InvitationCode) *adminv1.InvitationCode {
	return &adminv1.InvitationCode{
		Code:      ic.Code(),
		TeamCode:  int64(ic.Team().Code()),
		ExpiresAt: timestamppb.New(ic.ExpiresAt()),
	}
}
