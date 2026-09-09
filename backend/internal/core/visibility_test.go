package core

import (
	"testing"
	"time"
)

func TestMarkingVisibilityAtDelayFreezeAndFinalReveal(t *testing.T) {
	submitted := time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)
	freeze := submitted.Add(10 * time.Minute)
	final := submitted.Add(30 * time.Minute)
	tests := []struct {
		name   string
		now    time.Time
		freeze *time.Time
		final  *time.Time
		want   Visibility
	}{
		{name: "before delay", now: submitted.Add(PublishDelay - time.Nanosecond), want: VisibilityPrivate},
		{name: "delay boundary", now: submitted.Add(PublishDelay), want: VisibilityPublic},
		{name: "after freeze", now: submitted.Add(PublishDelay), freeze: &freeze, want: VisibilityTeam},
		{name: "final bypasses delay", now: submitted.Add(time.Minute), freeze: &freeze, final: &final, want: VisibilityPublic},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := MarkingVisibilityAt(submitted, test.now, test.freeze, test.final); got != test.want {
				t.Fatalf("visibility = %s, want %s", got, test.want)
			}
		})
	}
}
