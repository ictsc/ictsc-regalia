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
