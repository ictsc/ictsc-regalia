package domain

import (
	"context"
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	TeamRank struct {
		team           *Team
		rank           uint32
		totalScore     uint32
		updateSubmitAt *time.Time
	}
	Ranking []*TeamRank
)

type RankingEffect interface {
	TeamsLister
	RankingReader
}

func GetRankingForPublic(ctx context.Context, eff RankingEffect) (Ranking, error) {
	return getRanking(ctx, eff, true)
}

func GetRankingForAdmin(ctx context.Context, eff RankingEffect) (Ranking, error) {
	return getRanking(ctx, eff, false)
}

func getRanking(ctx context.Context, eff RankingEffect, isPublic bool) (Ranking, error) {
	teams, err := ListTeams(ctx, eff)
	if err != nil {
		return nil, err
	}

	ranks, err := eff.GetRanking(ctx, isPublic)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get ranking")
	}

	teamRanks := make([]*TeamRank, 0, len(teams))
	for _, team := range teams {
		teamRank := &TeamRank{team: team}

		idx := slices.IndexFunc(ranks, func(r *RankData) bool {
			return r.TeamID == uuid.UUID(team.teamID)
		})
		if idx >= 0 {
			teamRank.totalScore = ranks[idx].TotalScore
			teamRank.updateSubmitAt = &ranks[idx].UpdateSubmitAt
		}

		teamRanks = append(teamRanks, teamRank)
	}
	ranking := Ranking(teamRanks)

	ranking.sort()

	return ranking, nil
}

//nolint:cyclop
func (r Ranking) sort() {
	slices.SortStableFunc(r, func(lhs, rhs *TeamRank) int {
		// Descending order
		totalScoreCmp := rhs.totalScore - lhs.totalScore
		if totalScoreCmp != 0 {
			return int(totalScoreCmp)
		}

		// updateSubmitAt が nil のとき無限大として扱う
		switch {
		case lhs.updateSubmitAt == nil && rhs.updateSubmitAt != nil:
			return 1
		case lhs.updateSubmitAt != nil && rhs.updateSubmitAt == nil:
			return -1
		case lhs.updateSubmitAt == nil && rhs.updateSubmitAt == nil:
			return 0
		default:
			return lhs.updateSubmitAt.Compare(*rhs.updateSubmitAt)
		}
	})

	var (
		prev *TeamRank
		rank uint32 = 0
	)
	for _, teamRank := range r {
		if prev == nil ||
			prev.totalScore != teamRank.totalScore ||
			(prev.updateSubmitAt == nil) != (teamRank.updateSubmitAt == nil) ||
			(prev.updateSubmitAt != nil && *prev.updateSubmitAt != *teamRank.updateSubmitAt) {
			rank++
		}
		teamRank.rank = rank
		prev = teamRank
	}
}

func (t *TeamRank) Team() *Team {
	return t.team
}

func (t *TeamRank) Rank() uint32 {
	return t.rank
}

func (t *TeamRank) TotalScore() uint32 {
	return t.totalScore
}

func (t *TeamRank) UpdateSubmitAt() *time.Time {
	return t.updateSubmitAt
}

type (
	RankData struct {
		TeamID         uuid.UUID `json:"team_id"`
		TotalScore     uint32    `json:"total_score"`
		UpdateSubmitAt time.Time `json:"update_submit_at"`
	}
	RankingReader interface {
		GetRanking(ctx context.Context, isPublic bool) ([]*RankData, error)
	}
)
