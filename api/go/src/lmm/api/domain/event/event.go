package event

import (
	"time"

	"github.com/pkg/errors"
)

var ErrInvalidEvent = errors.New("invalid event")

type Event interface {
	Topic() string
	OccurredAt() time.Time
}
