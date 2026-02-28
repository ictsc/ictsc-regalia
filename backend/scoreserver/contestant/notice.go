package contestant

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type NoticeServiceHandler struct {
	contestantv1connect.UnimplementedNoticeServiceHandler

	ListEffect NoticeListEffect
}

var _ contestantv1connect.NoticeServiceHandler = (*NoticeServiceHandler)(nil)

func newNoticeServiceHandler(repo *pg.Repository) *NoticeServiceHandler {
	return &NoticeServiceHandler{
		ListEffect: repo,
	}
}

type NoticeListEffect interface {
	domain.NoticeReader
}

func (h *NoticeServiceHandler) ListNotices(
	ctx context.Context,
	req *connect.Request[contestantv1.ListNoticesRequest],
) (*connect.Response[contestantv1.ListNoticesResponse], error) {
	if _, err := session.UserSessionStore.Get(ctx); err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	notices, err := domain.ListEffectiveNotices(ctx, time.Now(), h.ListEffect)
	if err != nil {
		return nil, err
	}

	protoNotices := make([]*contestantv1.Notice, 0, len(notices))
	for _, notice := range notices {
		protoNotices = append(protoNotices, &contestantv1.Notice{
			Slug:  notice.Slug(),
			Title: notice.Title(),
			Body:  notice.Markdown(),
		})
	}

	return connect.NewResponse(&contestantv1.ListNoticesResponse{
		Notices: protoNotices,
	}), nil
}
