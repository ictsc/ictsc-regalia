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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ContestServiceHandler struct {
	contestantv1connect.UnimplementedContestServiceHandler

	GetRuleEffect     GetRuleEffect
	GetScheduleEffect domain.ScheduleReader
}

var _ contestantv1connect.ContestServiceHandler = (*ContestServiceHandler)(nil)

func newContestServiceHandler(repo *pg.Repository, scheduleReader domain.ScheduleReader) *ContestServiceHandler {
	return &ContestServiceHandler{
		GetRuleEffect:     repo,
		GetScheduleEffect: scheduleReader,
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

func (h *ContestServiceHandler) GetSchedule(
	ctx context.Context,
	req *connect.Request[contestantv1.GetScheduleRequest],
) (*connect.Response[contestantv1.GetScheduleResponse], error) {
	if _, err := session.UserSessionStore.Get(ctx); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	schedules, err := domain.GetSchedule(ctx, h.GetScheduleEffect)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	schedule := &contestantv1.Schedule{}

	// コンテストが開始済みか（いずれかのスケジュールが開始されたか）
	for _, entry := range schedules {
		if !now.Before(entry.StartAt()) {
			schedule.HasStarted = true
			break
		}
	}

	// 現在アクティブなスケジュール
	if current := schedules.Current(now); current != nil {
		schedule.Current = &contestantv1.ScheduleEntry{
			Name:    current.Name(),
			StartAt: timestamppb.New(current.StartAt()),
			EndAt:   timestamppb.New(current.EndAt()),
		}
	}

	// 次のスケジュール
	for _, entry := range schedules {
		if now.Before(entry.StartAt()) {
			schedule.Next = &contestantv1.ScheduleEntry{
				Name:    entry.Name(),
				StartAt: timestamppb.New(entry.StartAt()),
				EndAt:   timestamppb.New(entry.EndAt()),
			}
			break // schedules はstartAt順なので最初に見つかったものが次
		}
	}

	return connect.NewResponse(&contestantv1.GetScheduleResponse{
		Schedule: schedule,
	}), nil
}
