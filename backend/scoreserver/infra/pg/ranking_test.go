package pg_test

import (
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestGetRankingOrdering(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))
	ranking, err := repo.GetRanking(t.Context())
	if err != nil {
		t.Fatalf("failed to get ranking: %v", err)
	}

	// 順位順になっているかを確認:
	// ・前の要素のscoreが後の要素のscoreよりも小さくなっていない（降順）
	// ・scoreが等しい場合、前の要素のsubmitted_atが後の要素よりも早い（昇順）
	for i := 1; i < len(ranking); i++ {
		prev := ranking[i-1]
		curr := ranking[i]
		if prev.Score < curr.Score {
			t.Errorf("ranking is not sorted by score descending: index %d (score=%d) < index %d (score=%d)",
				i-1, prev.Score, i, curr.Score)
		}
		if prev.Score == curr.Score && prev.SubmittedAt.After(curr.SubmittedAt) {
			t.Errorf("ranking with equal score is not sorted by submitted_at ascending: index %d (submitted_at=%s) > index %d (submitted_at=%s) with score=%d",
				i-1, prev.SubmittedAt, i, curr.SubmittedAt, prev.Score)
		}
	}
}
