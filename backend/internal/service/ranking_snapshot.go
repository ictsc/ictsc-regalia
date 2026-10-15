package service

import (
	"context"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func (s *Service) frozenRanking(ctx context.Context, frozenAt time.Time) ([]core.RankingEntry, error) {
	if snapshot, ok, err := s.Store.GetRankingSnapshot(ctx, frozenAt); err != nil {
		return nil, err
	} else if ok {
		return snapshot.Entries, nil
	}
	content, err := s.Store.ContentAt(ctx, frozenAt)
	if err != nil {
		return nil, err
	}
	scores, err := s.publicScores(ctx, &frozenAt, false)
	if err != nil {
		return nil, err
	}
	teams, err := s.Store.ListTeams(ctx)
	if err != nil {
		return nil, err
	}
	active := make(map[string]struct{}, len(content.Manifest.Problems))
	for _, problem := range content.Manifest.Problems {
		active[problem.Code] = struct{}{}
	}
	stored, err := s.Store.PutRankingSnapshot(ctx, core.RankingSnapshot{
		FrozenAt: frozenAt, ContentCommit: content.CommitSHA, Entries: core.BuildRanking(teams, active, scores),
	}, "system:freeze-snapshot")
	if err != nil {
		return nil, err
	}
	return stored.Entries, nil
}
