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

type NoticeServiceHandler struct {
	Enforcer     *auth.Enforcer
	ListEffect   domain.NoticeReader
	UpdateEffect domain.Tx[domain.NoticeWriter]

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
		UpdateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.NoticeWriter { return rt }),
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
		protoNotices = append(protoNotices, &adminv1.Notice{
			Slug:          notice.Slug(),
			Title:         notice.Title(),
			Markdown:      notice.Markdown(),
			EffectiveFrom: timestamppb.New(notice.EffectiveFrom()),
		})
	}

	return connect.NewResponse(&adminv1.ListNoticesResponse{
		Notices: protoNotices,
	}), nil
}

func (s *NoticeServiceHandler) UpdateNotices(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateNoticesRequest],
) (*connect.Response[adminv1.UpdateNoticesResponse], error) {
	if err := enforce(ctx, s.Enforcer, "notices", "create"); err != nil {
		return nil, err
	}

	protoNotices := req.Msg.GetNotices()
	input := make(domain.NoticesInput, 0, len(protoNotices))
	for _, protoNotice := range protoNotices {
		input = append(input, &domain.NoticeData{
			Slug:          protoNotice.GetSlug(),
			Title:         protoNotice.GetTitle(),
			Markdown:      protoNotice.GetMarkdown(),
			EffectiveFrom: protoNotice.GetEffectiveFrom().AsTime(),
		})
	}
	notices, err := domain.NewNotices(input)
	if err != nil {
		return nil, err
	}

	if err := s.UpdateEffect.RunInTx(ctx, func(eff domain.NoticeWriter) error {
		return notices.Save(ctx, eff)
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateNoticesResponse{}), nil
}
