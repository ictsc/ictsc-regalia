package memory

import (
	"context"
	"net/http"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func (s *CompetitionStore) ContentAt(_ context.Context, at time.Time) (core.ContentSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var selected core.ContentSnapshot
	found := false
	for _, snapshot := range s.contents {
		if snapshot.ActivatedAt.After(at) {
			continue
		}
		if !found || snapshot.ActivatedAt.After(selected.ActivatedAt) {
			selected, found = snapshot, true
		}
	}
	if !found {
		return core.ContentSnapshot{}, core.NewError(http.StatusBadGateway, "content_not_available", "No valid content snapshot existed at the requested time")
	}
	return selected, nil
}

func (s *CompetitionStore) RecalculateScores(_ context.Context, selections map[int64]map[string]core.Score, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copySelections := make(map[int64]map[string]core.Score, len(selections))
	for team, byProblem := range selections {
		copySelections[team] = make(map[string]core.Score, len(byProblem))
		for problem, score := range byProblem {
			copySelections[team][problem] = score
		}
	}
	s.scoreSelections = copySelections
	return nil
}

func (s *CompetitionStore) GetRankingSnapshot(_ context.Context, frozenAt time.Time) (core.RankingSnapshot, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok := s.rankingSnapshots[frozenAt.UTC().Format(time.RFC3339Nano)]
	if ok {
		snapshot.Entries = append([]core.RankingEntry(nil), snapshot.Entries...)
	}
	return snapshot, ok, nil
}

func (s *CompetitionStore) PutRankingSnapshot(_ context.Context, snapshot core.RankingSnapshot, _ string) (core.RankingSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := snapshot.FrozenAt.UTC().Format(time.RFC3339Nano)
	if existing, ok := s.rankingSnapshots[key]; ok {
		return existing, nil
	}
	snapshot.Entries = append([]core.RankingEntry(nil), snapshot.Entries...)
	s.rankingSnapshots[key] = snapshot
	return snapshot, nil
}
