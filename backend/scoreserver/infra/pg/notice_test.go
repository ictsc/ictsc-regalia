package pg_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListNotices(t *testing.T) {
	t.Parallel()

	effectiveFrom, _ := time.Parse(time.RFC3339, "2025-02-03T00:00:00Z")
	effectiveUntil, _ := time.Parse(time.RFC3339, "2035-03-03T00:00:00Z")

	expected := []*domain.NoticeData{
		{
			ID:             uuid.FromStringOrNil("0cea0d50-96a5-45fb-a5c5-a6d6df140adc"),
			Path:           "/current/notice",
			Title:          "Current Notice",
			Markdown:       "現在のお知らせです",
			EffectiveFrom:  &effectiveFrom,
			EffectiveUntil: &effectiveUntil,
		},
		{
			ID:             uuid.FromStringOrNil("6ca38a12-adff-48f3-8fce-8f189eba38bb"),
			Path:           "/current/notice2",
			Title:          "Current Notice2",
			Markdown:       "現在のお知らせ2です",
			EffectiveFrom:  &effectiveFrom,
			EffectiveUntil: &effectiveUntil,
		},
	}

	repo := pg.NewRepository(pgtest.SetupDB(t))
	actual, err := repo.ListNotices(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	slices.SortFunc(actual, func(a, b *domain.NoticeData) int {
		return strings.Compare(a.ID.String(), b.ID.String())
	})

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("unexpected problems (-want +got):\n%s", diff)
	}
}

func TestSaveNotice(t *testing.T) {
	t.Parallel()

	effectiveFrom, _ := time.Parse(time.RFC3339, "2025-02-03T00:00:00Z")
	effectiveUntil, _ := time.Parse(time.RFC3339, "2035-03-03T00:00:00Z")
	effectiveUnitl2, _ := time.Parse(time.RFC3339, "2045-03-03T00:00:00Z")

	cases := map[string]struct {
		in    *domain.NoticeData
		query string
	}{
		"create": {
			in: &domain.NoticeData{
				ID:             uuid.FromStringOrNil("f3afb1e7-f281-4fbe-955e-1c91fef5c619"),
				Path:           "/test",
				Title:          "testお知らせです",
				Markdown:       "これはサンプルです。",
				EffectiveFrom:  &effectiveFrom,
				EffectiveUntil: &effectiveUntil,
			},
			query: `
			SELECT 1 FROM notices WHERE
			id = 'f3afb1e7-f281-4fbe-955e-1c91fef5c619' AND
			path = '/test' AND
			title = 'testお知らせです' AND
			markdown = 'これはサンプルです。' AND
			effective_from = '2025-02-03 00:00:00' AND
			effective_until = '2035-03-03 00:00:00'
			`,
		},
		"update": {
			in: &domain.NoticeData{
				ID:             uuid.FromStringOrNil("0cea0d50-96a5-45fb-a5c5-a6d6df140adc"),
				Path:           "/changed/notice",
				Title:          "Changed Notice",
				Markdown:       "変更済みのお知らせです",
				EffectiveFrom:  &effectiveUntil,
				EffectiveUntil: &effectiveUnitl2,
			},
			query: `
			SELECT 1 FROM notices WHERE
			id = '0cea0d50-96a5-45fb-a5c5-a6d6df140adc' AND
			path = '/changed/notice' AND
			title = 'Changed Notice' AND
			markdown = '変更済みのお知らせです' AND
			effective_from = '2035-03-03 00:00:00+00' AND
			effective_until = '2045-03-03 00:00:00'
			`,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := pgtest.SetupDB(t)
			repo := pg.NewRepository(db)

			if err := repo.RunTx(t.Context(), func(tx *pg.RepositoryTx) error {
				return tx.SaveNotice(t.Context(), tt.in)
			}); err != nil {
				t.Fatal(err)
			}

			var dst any
			if err := db.GetContext(t.Context(), &dst, tt.query); err != nil {
				t.Errorf("query %s: %v", tt.query, err)
			}
		})
	}
}
