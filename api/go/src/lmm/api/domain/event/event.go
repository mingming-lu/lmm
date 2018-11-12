package event

import "time"

type Event interface {
	Topic() string
	OccurredAt() time.Time
}
