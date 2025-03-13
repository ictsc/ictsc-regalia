package contestant

import (
	"context"

	"connectrpc.com/connect"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type RankingServiceHandler struct {
	GetEffect RankingGetEffect
}

var _ contestantv1connect.RankingServiceHandler = (*RankingServiceHandler)(nil)

func newRankingServiceHandler(repo *pg.Repository) *RankingServiceHandler {
	return &RankingServiceHandler{
		GetEffect: repo,
	}
}

type RankingGetEffect interface {
	domain.RankingGetter
}

func (r *RankingServiceHandler) GetRanking(
	ctx context.Context,
	req *connect.Request[contestantv1.GetRankingRequest],
) (*connect.Response[contestantv1.GetRankingResponse], error) {
	rankingData, err := domain.GetRanking(ctx, r.GetEffect)
	if err != nil {
		return nil, err
	}

	ranking := make([]*contestantv1.Rank, 0, len(rankingData))
	for _, rank := range rankingData {
		ranking = append(ranking, &contestantv1.Rank{
			Rank:     rank.Rank(),
			TeamName: rank.TeamName(),
			Score:    rank.Score(),
		})
	}

	return connect.NewResponse(&contestantv1.GetRankingResponse{
		Ranking: ranking,
	}), nil
}
