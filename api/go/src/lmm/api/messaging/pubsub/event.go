package pubsub

import (
	"time"

	"cloud.google.com/go/pubsub"
)

// Event wraps gcp pubsub message and implementes lmm/api/messaging.Evnet
type Event struct {
	topic string
	msg   *pubsub.Message
}

// NewEvent creates a new pubsub publish event
func NewEvent(topic string, data []byte) *Event {
	return &Event{
		topic: topic,
		msg:   &pubsub.Message{Data: data},
	}
}

// Topic implementation
func (e *Event) Topic() string {
	return e.topic
}

// OccurredAt implementation
func (e *Event) OccurredAt() time.Time {
	return e.msg.PublishTime
}
