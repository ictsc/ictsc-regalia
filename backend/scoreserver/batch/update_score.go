package batch

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
)

type ScoreUpdate struct {
	period     time.Duration
	markClient adminv1connect.MarkServiceClient
}

func NewScoreUpdate(
	cfg config.ScoreUpdate,
	markClient adminv1connect.MarkServiceClient,
) *ScoreUpdate {
	return &ScoreUpdate{
		period:     cfg.Period,
		markClient: markClient,
	}
}

func (s *ScoreUpdate) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Start update scores")

	ticker := time.NewTicker(s.period)
	defer ticker.Stop()

	for {
		if _, err := s.markClient.UpdateScores(
			ctx, connect.NewRequest(&adminv1.UpdateScoresRequest{}),
		); err != nil {
			slog.ErrorContext(ctx, "failed to update scores", "error", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}
