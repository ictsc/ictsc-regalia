package pg_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/pkg/snaptest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListNotices(t *testing.T) {
	t.Parallel()

	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.ListNotices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	snaptest.Match(t, actual, cmpopts.SortSlices(func(a, b *domain.NoticeData) int {
		return a.EffectiveFrom.Compare(b.EffectiveFrom)
	}))
}

func TestSaveNotices(t *testing.T) {
	t.Parallel()

	notices := []*domain.NoticeData{
		{
			Slug:          "test",
			Title:         "testお知らせです",
			Markdown:      "これはサンプルです。",
			EffectiveFrom: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Slug:          "changed/notice",
			Title:         "Changed Notice",
			Markdown:      "変更済みのお知らせです",
			EffectiveFrom: time.Date(2035, 3, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)

	if err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
		return tx.SaveNotices(t.Context(), notices)
	}); err != nil {
		t.Fatal(err)
	}

	actual, err := repo.ListNotices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	snaptest.Match(t, actual, cmpopts.SortSlices(func(a, b *domain.NoticeData) int {
		return a.EffectiveFrom.Compare(b.EffectiveFrom)
	}))
}
