package messaging

import (
	"time"
)

// Event interface
type Event interface {
	Topic() string
	Message() []byte
	PublishedAt() time.Time
}
