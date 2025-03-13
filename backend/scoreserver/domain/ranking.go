package domain

import (
	"context"
	"time"
)

type (
	Ranking struct {
		rank        int64
		teamName    string
		score       int64
		submittedAt time.Time
	}
	RankingData struct {
		TeamName    string    `db:"team_name"`
		Score       int64     `db:"score"`
		SubmittedAt time.Time `db:"submitted_at"`
	}
)

func (r *Ranking) Rank() int64 {
	return r.rank
}

func (r *Ranking) TeamName() string {
	return r.teamName
}

func (r *Ranking) Score() int64 {
	return r.score
}

func (r *Ranking) SubmittedAt() time.Time {
	return r.submittedAt
}

func (r *RankingData) parse() *Ranking {
	return &Ranking{
		teamName:    r.TeamName,
		score:       r.Score,
		submittedAt: r.SubmittedAt,
	}
}

func GetRanking(ctx context.Context, eff RankingGetter) ([]*Ranking, error) {
	rankingData, err := eff.GetRanking(ctx)
	if err != nil {
		return nil, err
	}
	ranking := make([]*Ranking, 0, len(rankingData))
	for _, rank := range rankingData {
		ranking = append(ranking, rank.parse())
	}
	AssignRanking(ranking)

	return ranking, nil
}

func AssignRanking(rankings []*Ranking) {
	var currentRank int64 = 0
	var prevScore int64
	var prevSubmittedAt time.Time

	for i, r := range rankings {
		if i == 0 || r.score != prevScore || !r.submittedAt.Equal(prevSubmittedAt) {
			currentRank = int64(i + 1)
		}
		r.rank = currentRank

		prevScore = r.score
		prevSubmittedAt = r.submittedAt
	}
}

type (
	RankingGetter interface {
		GetRanking(ctx context.Context) ([]*RankingData, error)
	}
)
