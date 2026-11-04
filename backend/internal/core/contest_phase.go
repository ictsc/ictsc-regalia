package core

import "time"

type ContestPhase string

const (
	// ContestPhaseBefore: no section is running and none has started yet.
	ContestPhaseBefore ContestPhase = "BEFORE"
	// ContestPhaseDuring: now is inside some section.
	ContestPhaseDuring ContestPhase = "DURING"
	// ContestPhaseAfter: no section is running but one has started, which
	// covers breaks between sections as well as the end of the contest.
	ContestPhaseAfter ContestPhase = "AFTER"
)

// devPhaseMargin is how far from now the shifted schedule starts or ends.
const devPhaseMargin = 10 * time.Minute

func ContestPhaseAt(sections []Section, now time.Time) ContestPhase {
	started := false
	for _, section := range sections {
		if now.Before(section.Beginning) {
			continue
		}
		if now.Before(section.Ending) {
			return ContestPhaseDuring
		}
		started = true
	}
	if started {
		return ContestPhaseAfter
	}
	return ContestPhaseBefore
}

// ShiftSectionsToPhase moves every section by the same offset so that now
// falls in phase, keeping each section's length and the breaks between them.
//   - BEFORE: the first section begins devPhaseMargin from now.
//   - DURING: the first section began devPhaseMargin ago (or halfway through
//     it, when it is shorter than twice the margin).
//   - AFTER: the last section ended devPhaseMargin ago.
func ShiftSectionsToPhase(sections []Section, phase ContestPhase, now time.Time) []Section {
	if len(sections) == 0 {
		return sections
	}
	first, last := sections[0], sections[0]
	for _, section := range sections[1:] {
		if section.Beginning.Before(first.Beginning) {
			first = section
		}
		if section.Ending.After(last.Ending) {
			last = section
		}
	}
	var offset time.Duration
	switch phase {
	case ContestPhaseBefore:
		offset = now.Add(devPhaseMargin).Sub(first.Beginning)
	case ContestPhaseDuring:
		lead := devPhaseMargin
		if half := first.Ending.Sub(first.Beginning) / 2; half < lead {
			lead = half
		}
		offset = now.Add(-lead).Sub(first.Beginning)
	case ContestPhaseAfter:
		offset = now.Add(-devPhaseMargin).Sub(last.Ending)
	}
	shifted := make([]Section, len(sections))
	for index, section := range sections {
		section.Beginning = section.Beginning.Add(offset).UTC().Truncate(time.Second)
		section.Ending = section.Ending.Add(offset).UTC().Truncate(time.Second)
		shifted[index] = section
	}
	return shifted
}
