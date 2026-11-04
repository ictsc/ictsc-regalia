package memory

import (
	"context"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func (s *CompetitionStore) ReplaceActiveSections(_ context.Context, sections []core.Section) (core.ContentSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	snapshot, ok := s.contents[s.active]
	if !ok {
		return core.ContentSnapshot{}, core.NewError(http.StatusBadGateway, "content_not_available", "No valid content snapshot is available")
	}
	snapshot.Manifest.Sections = append([]core.Section(nil), sections...)
	s.contents[s.active] = snapshot
	return snapshot, nil
}
