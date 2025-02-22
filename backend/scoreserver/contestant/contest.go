package contestant

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type ContestServiceHandler struct {
	contestantv1connect.UnimplementedContestServiceHandler

	GetRuleEffect GetRuleEffect
}

var _ contestantv1connect.ContestServiceHandler = (*ContestServiceHandler)(nil)

func newContestServiceHandler(repo *pg.Repository) *ContestServiceHandler {
	return &ContestServiceHandler{
		GetRuleEffect: repo,
	}
}

type GetRuleEffect interface {
	domain.RuleReader
}

func (h *ContestServiceHandler) GetRule(
	ctx context.Context,
	req *connect.Request[contestantv1.GetRuleRequest],
) (*connect.Response[contestantv1.GetRuleResponse], error) {
	if _, err := session.UserSessionStore.Get(ctx); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	rule, err := domain.GetRule(ctx, h.GetRuleEffect)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&contestantv1.GetRuleResponse{
		Rule: &contestantv1.Rule{
			Markdown: rule.Markdown(),
		},
	}), nil
}
