package admin

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/growi"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type NoticeServiceHandler struct {
	Enforcer     *auth.Enforcer
	ListEffect   NoticeListEffect
	CreateEffect NoticeCreateEffect

	adminv1connect.UnimplementedNoticeServiceHandler
}

var _ adminv1connect.NoticeServiceHandler = (*NoticeServiceHandler)(nil)

func NewNoticeServicehandler(
	enforcer *auth.Enforcer,
	repo *pg.Repository,
	growiClient *growi.Client,
) *NoticeServiceHandler {
	createEffect := struct {
		domain.NoticeWriter
		domain.NoticeGetter
	}{
		NoticeWriter: repo,
		NoticeGetter: growiClient,
	}
	return &NoticeServiceHandler{
		Enforcer:     enforcer,
		ListEffect:   repo,
		CreateEffect: createEffect,
	}
}

type NoticeListEffect = domain.NoticeReader
type NoticeCreateEffect interface {
	domain.NoticeWriter
	domain.NoticeGetter
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

func (s *NoticeServiceHandler) SyncNotices(
	ctx context.Context,
	req *connect.Request[adminv1.SyncNoticesRequest],
) (*connect.Response[adminv1.SyncNoticesResponse], error) {
	if err := enforce(ctx, s.Enforcer, "notices", "create"); err != nil {
		return nil, err
	}
	notice, err := domain.FetchNoticeByPath(ctx, s.CreateEffect, req.Msg.Path)
	if err != nil {
		return nil, err
	}

	if err := domain.SaveNotice(ctx, s.CreateEffect, notice); err != nil {
		return nil, err
	}
	protoNotice := convertNoticeData(notice)

	return connect.NewResponse(&adminv1.SyncNoticesResponse{
		Notice: protoNotice,
	}), nil
}

func convertNoticeData(notice *domain.NoticeData) *adminv1.Notice {
	return &adminv1.Notice{
		Path:     notice.Path,
		Title:    notice.Title,
		Markdown: notice.Markdown,
	}
}

func convertNotice(notice *domain.Notice) *adminv1.Notice {
	return &adminv1.Notice{
		Path:     notice.Path(),
		Title:    notice.Title(),
		Markdown: notice.Markdown(),
	}
}
