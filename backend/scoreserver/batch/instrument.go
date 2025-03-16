package batch

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func setupInstrumentation(client adminv1connect.RankingServiceClient) error {
	var errs []error

	if _, err := meter.Int64ObservableUpDownCounter(
		"ictsc.score",
		metric.WithDescription("Current score of each team and problem"),
		metric.WithUnit("1"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			scoresResp, err := client.ListScore(ctx, connect.NewRequest(&adminv1.ListScoreRequest{}))
			if err != nil {
				return errors.Wrap(err, "failed to list scores")
			}
			for _, score := range scoresResp.Msg.GetScores() {
				observer.Observe(
					score.GetScore(),
					metric.WithAttributes(
						attribute.Int64("ictsc.team.code", score.GetTeam().GetCode()),
						attribute.String("ictsc.team.name", score.GetTeam().GetName()),
						attribute.String("ictsc.problem.code", score.GetProblem().GetCode()),
						attribute.String("ictsc.problem.category", score.GetProblem().GetCategory()),
					),
				)
			}
			return nil
		}),
	); err != nil {
		errs = append(errs, errors.Wrap(err, "failed to create ictsc.score gauge"))
	}

	return errors.Join(errs...)
}
