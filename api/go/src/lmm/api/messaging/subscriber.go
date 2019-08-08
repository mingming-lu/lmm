package messaging

import "context"

// EventHandler type
type EventHandler func(context.Context, Event) error

// Subscriber interface
type Subscriber interface {
	Subscribe(e Event, handler EventHandler) error
}
