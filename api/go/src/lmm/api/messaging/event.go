package messaging

import (
	"time"
)

// Event interface
type Event interface {
	Topic() string
	Message() interface{}
	PublishedAt() time.Time
}
