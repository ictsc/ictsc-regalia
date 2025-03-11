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

func newContestServiceHandler(repo *pg.Repository) *ContestServiceHandler {
	return &ContestServiceHandler{
		GetRuleEffect:     repo,
		GetScheduleEffect: repo,
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

	schedule, err := domain.GetSchedule(ctx, h.GetScheduleEffect)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	current := schedule.Current(now)
	next := schedule.Next(now)

	protoSchedule := &contestantv1.Schedule{
		Phase:     convertPhase(current.Phase()),
		NextPhase: convertPhase(next.Phase()),
		StartAt:   timestamppb.New(current.StartAt()),
	}
	if endAt := current.EndAt(); !endAt.IsZero() {
		protoSchedule.EndAt = timestamppb.New(endAt)
	}

	return connect.NewResponse(&contestantv1.GetScheduleResponse{
		Schedule: protoSchedule,
	}), nil
}

func convertPhase(phase domain.Phase) contestantv1.Phase {
	switch phase {
	case domain.PhaseOutOfContest:
		return contestantv1.Phase_PHASE_OUT_OF_CONTEST
	case domain.PhaseInContest:
		return contestantv1.Phase_PHASE_IN_CONTEST
	case domain.PhaseBreak:
		return contestantv1.Phase_PHASE_BREAK
	case domain.PhaseAfterContest:
		return contestantv1.Phase_PHASE_AFTER_CONTEST
	case domain.PhaseUnspecified:
		fallthrough
	default:
		return contestantv1.Phase_PHASE_UNSPECIFIED
	}
}
