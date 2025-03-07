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

type ScheduleServiceHandler struct {
	adminv1connect.UnimplementedScheduleServiceHandler

	Enforcer     *auth.Enforcer
	GetEffect    ScheduleGetEffect
	UpdateEffect ScheduleUpdateEffect
}

var _ adminv1connect.ScheduleServiceHandler = (*ScheduleServiceHandler)(nil)

func newScheduleServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *ScheduleServiceHandler {
	return &ScheduleServiceHandler{
		Enforcer:     enforcer,
		GetEffect:    repo,
		UpdateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.ScheduleWriter { return rt }),
	}
}

type ScheduleGetEffect = domain.ScheduleReader
type ScheduleUpdateEffect = domain.Tx[domain.ScheduleWriter]

func (h *ScheduleServiceHandler) GetSchedule(
	ctx context.Context,
	req *connect.Request[adminv1.GetScheduleRequest],
) (*connect.Response[adminv1.GetScheduleResponse], error) {
	if err := enforce(ctx, h.Enforcer, "schedules", "get"); err != nil {
		return nil, err
	}

	schedules, err := domain.GetSchedule(ctx, h.GetEffect)
	if err != nil {
		return nil, err
	}

	protoSchedules := make([]*adminv1.Schedule, 0, len(schedules))
	for _, schedule := range schedules {
		protoSchedules = append(protoSchedules, convertSchedule(schedule))
	}

	return connect.NewResponse(&adminv1.GetScheduleResponse{
		Schedule: protoSchedules,
	}), nil
}

func (h *ScheduleServiceHandler) UpdateSchedule(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateScheduleRequest],
) (*connect.Response[adminv1.UpdateScheduleResponse], error) {
	if err := enforce(ctx, h.Enforcer, "schedules", "update"); err != nil {
		return nil, err
	}

	schedules := make([]*domain.UpdateScheduleInput, 0, len(req.Msg.GetSchedule()))
	for _, schdule := range req.Msg.GetSchedule() {
		schedules = append(schedules, &domain.UpdateScheduleInput{
			Phase:   convertProtoPhaseToDomain(schdule.GetPhase()),
			StartAt: schdule.GetStartAt().AsTime(),
			EndAt:   schdule.GetEndAt().AsTime(),
		})
	}
	if err := h.UpdateEffect.RunInTx(ctx, func(w domain.ScheduleWriter) error {
		return domain.SaveSchedule(ctx, w, schedules)
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateScheduleResponse{}), nil
}

func convertSchedule(schedule *domain.Schedule) *adminv1.Schedule {
	var phase adminv1.Phase
	switch schedule.Phase() {
	case domain.PhaseOutOfContest:
		phase = adminv1.Phase_PHASE_OUT_OF_CONTEST
	case domain.PhaseInContest:
		phase = adminv1.Phase_PHASE_IN_CONTEST
	case domain.PhaseBreak:
		phase = adminv1.Phase_PHASE_BREAK
	case domain.PhaseAfterContest:
		phase = adminv1.Phase_PHASE_AFTER_CONTEST
	case domain.PhaseUnspecified:
		fallthrough
	default:
		phase = adminv1.Phase_PHASE_UNSPECIFIED
	}

	return &adminv1.Schedule{
		Phase:   phase,
		StartAt: timestamppb.New(schedule.StartAt()),
		EndAt:   timestamppb.New(schedule.EndAt()),
	}
}

func convertProtoPhaseToDomain(protoPhase adminv1.Phase) domain.Phase {
	switch protoPhase {
	case adminv1.Phase_PHASE_OUT_OF_CONTEST:
		return domain.PhaseOutOfContest
	case adminv1.Phase_PHASE_IN_CONTEST:
		return domain.PhaseInContest
	case adminv1.Phase_PHASE_BREAK:
		return domain.PhaseBreak
	case adminv1.Phase_PHASE_AFTER_CONTEST:
		return domain.PhaseAfterContest
	default:
		return domain.PhaseUnspecified
	}
}
