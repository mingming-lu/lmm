package testing

import (
	"time"

	"lmm/api/clock"
)

type testClock struct{}

func (c testClock) Now() time.Time {
	return time.Now().Truncate(1 * time.Second)
}

func init() {
	clock.DefaultClock = &testClock{}
}
