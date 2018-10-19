package messaging

import "time"

// Event definition
type Event interface {
	Topic() string
	OccurredAt() time.Time
}
