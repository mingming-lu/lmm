package clock

import "time"

var DefaultClock Clock

type Clock interface {
	Now() time.Time
}

type clock struct{}

func (c clock) Now() time.Time {
	return time.Now().Truncate(1 * time.Second).UTC()
}

func Now() time.Time {
	return DefaultClock.Now()
}

func init() {
	DefaultClock = &clock{}
}
