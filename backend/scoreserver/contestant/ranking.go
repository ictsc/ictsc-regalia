package contestant

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RankingServiceHandler struct {
	contestantv1connect.UnimplementedRankingServiceHandler

	GetEffect domain.RankingEffect
}

var _ contestantv1connect.RankingServiceHandler = (*RankingServiceHandler)(nil)

func newRankingServiceHandler(repo domain.RankingEffect) *RankingServiceHandler {
	return &RankingServiceHandler{
		GetEffect: repo,
	}
}

func (h *RankingServiceHandler) GetRanking(
	ctx context.Context,
	req *connect.Request[contestantv1.GetRankingRequest],
) (*connect.Response[contestantv1.GetRankingResponse], error) {
	if _, err := session.UserSessionStore.Get(ctx); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}

	ranking, err := domain.GetRankingForPublic(ctx, h.GetEffect)
	if err != nil {
		return nil, err
	}

	protoRanks := make([]*contestantv1.Rank, 0, len(ranking))
	for _, rank := range ranking {
		proto := &contestantv1.Rank{
			Rank:         int64(rank.Rank()),
			TeamName:     rank.Team().Name(),
			Organization: rank.Team().Organization(),
			Score:        int64(rank.TotalScore()),
		}
		if lastEffectiveSubmitAt := rank.UpdateSubmitAt(); lastEffectiveSubmitAt != nil {
			proto.Timestamp = timestamppb.New(*lastEffectiveSubmitAt)
		}
		protoRanks = append(protoRanks, proto)
	}

	return connect.NewResponse(&contestantv1.GetRankingResponse{
		Ranking: protoRanks,
	}), nil
}
