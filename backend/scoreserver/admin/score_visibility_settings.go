package admin

import (
	"context"
	"time"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ScoreVisibilitySettingsServiceHandler struct {
	adminv1connect.UnimplementedScoreVisibilitySettingsServiceHandler

	Enforcer     *auth.Enforcer
	GetEffect    domain.ScoreVisibilitySettingsReader
	UpdateEffect domain.Tx[domain.ScoreVisibilitySettingsWriter]
}

var _ adminv1connect.ScoreVisibilitySettingsServiceHandler = (*ScoreVisibilitySettingsServiceHandler)(nil)

func newScoreVisibilitySettingsServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *ScoreVisibilitySettingsServiceHandler {
	return &ScoreVisibilitySettingsServiceHandler{
		Enforcer:  enforcer,
		GetEffect: repo,
		UpdateEffect: pg.Tx(repo, func(rt *pg.RepositoryTx) domain.ScoreVisibilitySettingsWriter {
			return rt
		}),
	}
}

func (h *ScoreVisibilitySettingsServiceHandler) GetScoreVisibilitySettings(
	ctx context.Context,
	req *connect.Request[adminv1.GetScoreVisibilitySettingsRequest],
) (*connect.Response[adminv1.GetScoreVisibilitySettingsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "score_visibility_settings", "get"); err != nil {
		return nil, err
	}

	settings, err := domain.GetScoreVisibilitySettings(ctx, h.GetEffect)
	if err != nil {
		return nil, err
	}

	resp := &adminv1.GetScoreVisibilitySettingsResponse{
		Settings: &adminv1.ScoreVisibilitySettings{},
	}
	if freezeAt := settings.RankingFreezeAt(); freezeAt != nil {
		resp.Settings.RankingFreezeAt = timestamppb.New(*freezeAt)
	}
	return connect.NewResponse(resp), nil
}

func (h *ScoreVisibilitySettingsServiceHandler) UpdateScoreVisibilitySettings(
	ctx context.Context,
	req *connect.Request[adminv1.UpdateScoreVisibilitySettingsRequest],
) (*connect.Response[adminv1.UpdateScoreVisibilitySettingsResponse], error) {
	if err := enforce(ctx, h.Enforcer, "score_visibility_settings", "update"); err != nil {
		return nil, err
	}

	var freezeAt *time.Time
	if ts := req.Msg.GetSettings().GetRankingFreezeAt(); ts != nil {
		t := ts.AsTime()
		freezeAt = &t
	}

	if err := h.UpdateEffect.RunInTx(ctx, func(w domain.ScoreVisibilitySettingsWriter) error {
		return domain.SaveScoreVisibilitySettings(ctx, w, &domain.UpdateScoreVisibilitySettingsInput{
			RankingFreezeAt: freezeAt,
		})
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&adminv1.UpdateScoreVisibilitySettingsResponse{}), nil
}
