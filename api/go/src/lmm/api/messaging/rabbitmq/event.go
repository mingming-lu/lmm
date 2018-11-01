package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
)

// Event defines rabbitmq event
type Event struct {
	msg *amqp.Publishing
}

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
