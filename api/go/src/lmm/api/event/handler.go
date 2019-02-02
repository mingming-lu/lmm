package event

import "context"

type EventHandler = func(context.Context, Event) error

var NopEventHandler = func(context.Context, Event) error {
	return nil
}
