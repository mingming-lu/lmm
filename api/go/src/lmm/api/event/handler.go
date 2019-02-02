package event

import "context"

type EventHandler = func(context.Context, Event) error
