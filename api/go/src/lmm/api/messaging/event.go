package messaging

import (
	"time"

	"github.com/pkg/errors"
)

var ErrInvalidEvent = errors.New("invalid event")

// Event definition
type Event interface {
	Topic() string
	OccurredAt() time.Time
}
