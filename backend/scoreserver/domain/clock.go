package domain

import "time"

type Clock func() time.Time

func (c Clock) Now() time.Time {
	if c == nil {
		return time.Now()
	}
	return c()
}
