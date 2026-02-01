package admin

import (
	"context"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RankingServiceHandler struct {
	adminv1connect.UnimplementedRankingServiceHandler

	Enforcer         *auth.Enforcer
	ListScoreEffect  domain.TeamProblemLister
	GetRankingEffect domain.RankingEffect
	ScheduleResolver ScheduleIDResolver
}

var _ adminv1connect.RankingServiceHandler = (*RankingServiceHandler)(nil)

func newRankingServiceHandler(enforcer *auth.Enforcer, repo *pg.Repository) *RankingServiceHandler {
	return &RankingServiceHandler{
		Enforcer:         enforcer,
		ListScoreEffect:  repo,
		GetRankingEffect: repo,
		ScheduleResolver: repo,
	}
}

func (h *RankingServiceHandler) ListScore(
	ctx context.Context,
	req *connect.Request[adminv1.ListScoreRequest],
) (*connect.Response[adminv1.ListScoreResponse], error) {
	if err := enforce(ctx, h.Enforcer, "scores", "list"); err != nil {
		return nil, err
	}

	teamProblems, err := domain.ListTeamProblemsForAdmin(ctx, h.ListScoreEffect)
	if err != nil {
		return nil, err
	}

	scheduleIDsSet := make(map[uuid.UUID]struct{})
	for _, teamProblem := range teamProblems {
		for _, scheduleID := range teamProblem.Problem().SubmissionableScheduleIDs() {
			scheduleIDsSet[scheduleID] = struct{}{}
		}
	}

	scheduleIDs := make([]uuid.UUID, 0, len(scheduleIDsSet))
	for id := range scheduleIDsSet {
		scheduleIDs = append(scheduleIDs, id)
	}

	scheduleNames, err := h.ScheduleResolver.GetScheduleNamesByIDs(ctx, scheduleIDs)
	if err != nil {
		return nil, err
	}

	protoScores := make([]*adminv1.Score, 0, len(teamProblems))
	for _, teamProblem := range teamProblems {
		score := teamProblem.Score()
		if score == nil {
			continue
		}
		protoScores = append(protoScores, &adminv1.Score{
			Team:        convertTeam(teamProblem.Team()),
			Problem:     convertProblem(teamProblem.Problem(), scheduleNames),
			MarkedScore: int64(score.MarkedScore()),
			Penalty:     int64(score.Penalty()),
			Score:       int64(score.TotalScore()),
		})
	}

	return connect.NewResponse(&adminv1.ListScoreResponse{
		Scores: protoScores,
	}), nil
}

func (h *RankingServiceHandler) GetRanking(
	ctx context.Context,
	req *connect.Request[adminv1.GetRankingRequest],
) (*connect.Response[adminv1.GetRankingResponse], error) {
	if err := enforce(ctx, h.Enforcer, "scores", "list"); err != nil {
		return nil, err
	}

	ranking, err := domain.GetRankingForAdmin(ctx, h.GetRankingEffect)
	if err != nil {
		return nil, err
	}

	protoRanks := make([]*adminv1.TeamRank, 0, len(ranking))
	for _, rank := range ranking {
		protoRank := &adminv1.TeamRank{
			Team:  convertTeam(rank.Team()),
			Rank:  int64(rank.Rank()),
			Score: int64(rank.TotalScore()),
		}
		if lastSubmitAt := rank.UpdateSubmitAt(); lastSubmitAt != nil {
			protoRank.LastEffectiveSubmissionAt = timestamppb.New(*lastSubmitAt)
		}
		protoRanks = append(protoRanks, protoRank)
	}

	return connect.NewResponse(&adminv1.GetRankingResponse{
		Ranking: protoRanks,
	}), nil
}
