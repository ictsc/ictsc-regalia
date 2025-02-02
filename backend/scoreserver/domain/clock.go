package domain

import "time"

type Clocker interface {
	Now() time.Time
}

type ClockerFunc func() time.Time

func (f ClockerFunc) Now() time.Time {
	return f()
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}
