package service

import (
	"context"
	"net/http"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

// DevContestPhaseEnabled reports whether admins may move the schedule to try
// the contest phases: development fake mode with secondary qualifier mode.
func (s *Service) DevContestPhaseEnabled() bool {
	return s.Config.DevTools && s.Config.SecondaryQualifierMode
}

// ContestPhase returns the phase of the active schedule at the current time.
func (s *Service) ContestPhase(ctx context.Context) (core.ContestPhase, []core.Section, error) {
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		return "", nil, err
	}
	return core.ContestPhaseAt(snapshot.Manifest.Sections, s.Now()), snapshot.Manifest.Sections, nil
}

// SetDevContestPhase shifts the active schedule relative to now so the
// contest is in phase. It rewrites the stored manifest and is refused outside
// development secondary qualifier mode.
func (s *Service) SetDevContestPhase(ctx context.Context, phase core.ContestPhase) (core.ContestPhase, []core.Section, error) {
	if !s.DevContestPhaseEnabled() {
		return "", nil, core.NewError(http.StatusNotFound, "resource_not_found", "Contest phase switching is only available in development secondary qualifier mode")
	}
	snapshot, err := s.Store.ActiveContent(ctx)
	if err != nil {
		return "", nil, err
	}
	shifted := core.ShiftSectionsToPhase(snapshot.Manifest.Sections, phase, s.Now())
	if _, err := s.Store.ReplaceActiveSections(ctx, shifted); err != nil {
		return "", nil, err
	}
	return s.ContestPhase(ctx)
}
