package clock

import "time"

var defaultClock clock

type Clock interface {
	Now() time.Time
}

type clock func() time.Time

func (c clock) Now() time.Time {
	return c().Truncate(1 * time.Second)
}

func Now() time.Time {
	return defaultClock.Now().UTC()
}

func init() {
	defaultClock = time.Now
}
