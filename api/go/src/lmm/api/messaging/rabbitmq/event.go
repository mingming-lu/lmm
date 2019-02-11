package rabbitmq

import (
	"lmm/api/messaging"
	"time"

	"github.com/streadway/amqp"
)

// Event defines rabbitmq event
type Event struct {
	messaging.Event

	msg *amqp.Publishing
}

// NewEvent creates a new rabbitmq event
func NewEvent(msg *amqp.Publishing) *Event {
	return &Event{
		msg: msg,
	}
}

// OccurredAt implementation
func (e *Event) OccurredAt() time.Time {
	return e.msg.Timestamp
}

// Topic implementation
func (e *Event) Topic() string {
	return e.msg.Type
}

// Message returns event's message
func (e *Event) Message() amqp.Publishing {
	return *e.msg
}
