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

type RuleServiceHandler struct {
	adminv1connect.UnimplementedRuleServiceHandler

	Enforcer     *auth.Enforcer
	GetEffect    RuleGetEffect
	UpdateEffect RuleUpdateEffect
}

var _ adminv1connect.RuleServiceHandler = (*RuleServiceHandler)(nil)

func newRuleServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *RuleServiceHandler {
	return &RuleServiceHandler{
		Enforcer:     enforcer,
		GetEffect:    repo,
		UpdateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.RuleWriter { return rt }),
	}
}

type RuleGetEffect interface {
	domain.RuleReader
}

func (h *RuleServiceHandler) GetRule(
	ctx context.Context,
	req *connect.Request[adminv1.GetRuleRequest],
) (*connect.Response[adminv1.GetRuleResponse], error) {
	if err := enforce(ctx, h.Enforcer, "rule", "get"); err != nil {
		return nil, err
	}

	rule, err := domain.GetRule(ctx, h.GetEffect)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetRuleResponse{
		Rule: &adminv1.Rule{
			Markdown: rule.Markdown(),
		},
	}), nil
}

type RuleUpdateEffect interface {
	domain.Tx[domain.RuleWriter]
}

func (h *RuleServiceHandler) UpdateRule(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateRuleRequest],
) (*connect.Response[adminv1.UpdateRuleResponse], error) {
	if err := enforce(ctx, h.Enforcer, "rule", "update"); err != nil {
		return nil, err
	}

	rule, err := domain.NewRule(req.Msg.GetRule().GetMarkdown())
	if err != nil {
		return nil, err
	}

	if err := h.UpdateEffect.RunInTx(ctx, func(w domain.RuleWriter) error {
		return rule.Save(ctx, w)
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateRuleResponse{}), nil
}
