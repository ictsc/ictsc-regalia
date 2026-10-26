package core

import (
	"regexp"
	"time"
)

// imageReference matches the /images/<file> URLs that content Markdown uses.
var imageReference = regexp.MustCompile(`/images/([A-Za-z0-9._-]+)`)

// ContestantVisibleImages lists the images contestants may load at now: the
// ones referenced by problems whose section has started, by announcements in
// effect, and by the rule. Explanations are never included.
func ContestantVisibleImages(snapshot ContentSnapshot, now time.Time) map[string]struct{} {
	started := make(map[string]bool, len(snapshot.Manifest.Sections))
	for _, section := range snapshot.Manifest.Sections {
		if !now.Before(section.Beginning) {
			started[section.Slug] = true
		}
	}
	visible := make(map[string]struct{})
	collect := func(markdown string) {
		for _, match := range imageReference.FindAllStringSubmatch(markdown, -1) {
			visible[match[1]] = struct{}{}
		}
	}
	for _, problem := range snapshot.Manifest.Problems {
		if started[problem.SectionSlug] {
			collect(problem.Body)
		}
	}
	for _, announcement := range snapshot.Manifest.Announcements {
		if !now.Before(announcement.EffectiveFrom) {
			collect(announcement.Markdown)
		}
	}
	collect(snapshot.Manifest.RuleMarkdown)
	return visible
}
