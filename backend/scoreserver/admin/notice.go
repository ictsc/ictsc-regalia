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

type NoticeServiceHandler struct {
	Enforcer     *auth.Enforcer
	ListEffect   domain.NoticeReader
	CreateEffect domain.Tx[domain.NoticeWriter]

	adminv1connect.UnimplementedNoticeServiceHandler
}

var _ adminv1connect.NoticeServiceHandler = (*NoticeServiceHandler)(nil)

func NewNoticeServicehandler(
	enforcer *auth.Enforcer,
	repo *pg.Repository,
) *NoticeServiceHandler {
	return &NoticeServiceHandler{
		Enforcer:     enforcer,
		ListEffect:   repo,
		CreateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.NoticeWriter { return rt }),
	}
}

func (s *NoticeServiceHandler) ListNotices(
	ctx context.Context,
	req *connect.Request[adminv1.ListNoticesRequest],
) (*connect.Response[adminv1.ListNoticesResponse], error) {
	if err := enforce(ctx, s.Enforcer, "notices", "list"); err != nil {
		return nil, err
	}
	notices, err := domain.ListNotices(ctx, s.ListEffect)
	if err != nil {
		return nil, err
	}

	protoNotices := make([]*adminv1.Notice, 0, len(notices))
	for _, notice := range notices {
		protoNotices = append(protoNotices, convertNotice(notice))
	}

	return connect.NewResponse(&adminv1.ListNoticesResponse{
		Notices: protoNotices,
	}), nil
}

func convertNotice(notice *domain.Notice) *adminv1.Notice {
	return &adminv1.Notice{
		Path:     notice.Path(),
		Title:    notice.Title(),
		Markdown: notice.Markdown(),
	}
}
