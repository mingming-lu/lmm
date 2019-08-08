package messaging

import "context"

// Publisher interface
type Publisher interface {
	Publish(context.Context, Event) error
}
