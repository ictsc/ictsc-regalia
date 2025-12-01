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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvitationServiceHandler struct {
	adminv1connect.UnimplementedInvitationServiceHandler
	Enforcer     *auth.Enforcer
	ListEffect   domain.InvitationCodeLister
	CreateEffect createInvitationCodeEffect
}

func NewInvitationServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *InvitationServiceHandler {
	return &InvitationServiceHandler{
		UnimplementedInvitationServiceHandler: adminv1connect.UnimplementedInvitationServiceHandler{},

		Enforcer:   enforcer,
		ListEffect: repo,
		CreateEffect: struct {
			*pg.Repository
			domain.SystemClock
		}{Repository: repo},
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

	ics, err := domain.ListInvitationCodes(ctx, h.ListEffect)
	if err != nil {
		return nil, err
	}

	protoICs := make([]*adminv1.InvitationCode, 0, len(ics))
	for _, ic := range ics {
		protoICs = append(protoICs, convertInvitationCode(ic))
	}

	return connect.NewResponse(&adminv1.ListInvitationCodesResponse{
		InvitationCodes: protoICs,
	}), nil
}

type createInvitationCodeEffect interface {
	domain.TeamGetter
	domain.InvitationCodeCreator
	domain.Clocker
}

func (h *InvitationServiceHandler) CreateInvitationCode(
	ctx context.Context,
	req *connect.Request[adminv1.CreateInvitationCodeRequest],
) (*connect.Response[adminv1.CreateInvitationCodeResponse], error) {
	if err := enforce(ctx, h.Enforcer, "invitation_codes", "create"); err != nil {
		return nil, err
	}

	now := h.CreateEffect.Now()

	protoIC := req.Msg.GetInvitationCode()
	protoTeamCode := protoIC.GetTeamCode()
	if protoTeamCode == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("team_code is required"))
	}
	teamCode, err := domain.NewTeamCode(protoTeamCode)
	if err != nil {
		return nil, err
	}

	team, err := teamCode.Team(ctx, h.CreateEffect)
	if err != nil {
		return nil, err
	}

	// リクエストからコードを取得（空文字列または不適切な文字列の場合は自動生成される）
	manualCode := protoIC.GetCode()

	invitationCode, err := team.CreateInvitationCode(ctx, h.CreateEffect, now, protoIC.GetExpiresAt().AsTime(), manualCode)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.CreateInvitationCodeResponse{
		InvitationCode: convertInvitationCode(invitationCode),
	}), nil
}

func convertInvitationCode(ic *domain.InvitationCode) *adminv1.InvitationCode {
	return &adminv1.InvitationCode{
		Code:      ic.Code(),
		TeamCode:  int64(ic.Team().Code()),
		ExpiresAt: timestamppb.New(ic.ExpiresAt()),
	}
}
