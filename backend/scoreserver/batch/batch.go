package batch

import (
	"context"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/sstate"
	"golang.org/x/oauth2"
)

type Batch struct {
	deploymentSync *DeploymentSync
	scoreUpdate    *ScoreUpdate
}

func NewBatch(cfg config.Batch) (*Batch, error) {
	apiClient := http.DefaultClient
	if cfg.APITokenSource != nil {
		apiClient.Transport = &oauth2.Transport{Source: cfg.APITokenSource, Base: apiClient.Transport}
	}

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OpenTelemetry interceptor")
	}
	connectOpts := []connect.ClientOption{
		connect.WithInterceptors(otelInterceptor),
		connect.WithGRPC(),
	}

	deploymentClient := adminv1connect.NewDeploymentServiceClient(apiClient, cfg.APIURL, connectOpts...)
	markClient := adminv1connect.NewMarkServiceClient(apiClient, cfg.APIURL, connectOpts...)
	rankingClient := adminv1connect.NewRankingServiceClient(apiClient, cfg.APIURL, connectOpts...)

	sstateClient, err := sstate.NewSStateClient(cfg.DeploymentSync.SState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sstate client")
	}

	var deploymentSync *DeploymentSync
	if cfg.DeploymentSync != nil {
		deploymentSync = NewDeploymentSync(*cfg.DeploymentSync, deploymentClient, sstateClient)
	}

	var scoreUpdate *ScoreUpdate
	if cfg.ScoreUpdate != nil {
		scoreUpdate = NewScoreUpdate(*cfg.ScoreUpdate, markClient)
	}

	if err := setupInstrumentation(rankingClient); err != nil {
		return nil, errors.Wrap(err, "failed to setup instrumentation")
	}

	return &Batch{
		deploymentSync: deploymentSync,
		scoreUpdate:    scoreUpdate,
	}, nil
}

func (b *Batch) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	if b.deploymentSync != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := b.deploymentSync.Run(ctx); err != nil {
				cancel(errors.Wrap(err, "failed to start sync deployments"))
			}
		}()
	}

	if b.scoreUpdate != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := b.scoreUpdate.Run(ctx); err != nil {
				cancel(errors.Wrap(err, "failed to start update scores"))
			}
		}()
	}

	wg.Wait()
	if errors.Is(ctx.Err(), context.Canceled) {
		if err := context.Cause(ctx); !errors.Is(err, context.Canceled) {
			return errors.WithStack(err)
		}
	}
	return nil
}
