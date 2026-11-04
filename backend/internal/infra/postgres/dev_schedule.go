package postgres

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func (s *Store) ReplaceActiveSections(ctx context.Context, sections []core.Section) (core.ContentSnapshot, error) {
	encoded, err := json.Marshal(sections)
	if err != nil {
		return core.ContentSnapshot{}, core.WrapError(http.StatusInternalServerError, "internal_error", "Sections cannot be encoded", err)
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE content_snapshots SET manifest=jsonb_set(manifest,'{sections}',$1::jsonb)
		WHERE active=true`, encoded)
	if err != nil {
		return core.ContentSnapshot{}, dbError(err)
	}
	if tag.RowsAffected() == 0 {
		return core.ContentSnapshot{}, core.NewError(http.StatusBadGateway, "content_not_available", "No valid content snapshot is available")
	}
	return s.ActiveContent(ctx)
}
