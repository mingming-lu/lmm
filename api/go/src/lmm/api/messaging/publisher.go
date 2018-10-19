package messaging

import "context"

// Publisher publishes events
type Publisher interface {
	Publish(context.Context, Event) error
}
