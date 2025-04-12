package tinylfu

import "time"

type Clock interface {
	Now() time.Time
}

type nativeClock struct{}

func (nativeClock) Now() time.Time {
	return time.Now()
}
