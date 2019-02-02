package event

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

var (
	ErrEmptyEventTopic   = errors.New("event topic cannot be empty")
	ErrEventHandleFailed = errors.New("failed to handle event")
	ErrUnRegisteredEvent = errors.New("unregistered event")
)

var (
	globalBusLock sync.RWMutex
	busInstance   Bus
)

func init() {
	busInstance = newBus()
}

// Bus is a pub/sub bus
type Bus interface {
	Publisher
	Subscriber
}

type bus struct {
	handlers map[string][]EventHandler
}

func (b *bus) Publish(c context.Context, e Event) error {
	globalBusLock.RLock()

	if e.Topic() == "" {
		return ErrEmptyEventTopic
	}

	handlers, ok := b.handlers[e.Topic()]
	globalBusLock.RUnlock()

	if !ok || len(handlers) == 0 {
		return errors.Wrap(ErrUnRegisteredEvent, e.Topic())
	}

	group, ctx := errgroup.WithContext(c)

	for _, handler := range handlers {
		group.Go(func() error {
			return handler(ctx, e)
		})
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- group.Wait()
	}()

	select {
	case err := <-errChan:
		return err
	case <-c.Done():
		return c.Err()
	}
}

func (b *bus) Subscribe(e Event, handler EventHandler) error {
	globalBusLock.Lock()
	defer globalBusLock.Unlock()

	if e.Topic() == "" {
		return ErrEmptyEventTopic
	}

	b.handlers[e.Topic()] = append(b.handlers[e.Topic()], handler)

	return nil
}

func newBus() Bus {
	return &bus{
		handlers: make(map[string][]EventHandler),
	}
}

// SyncBus gets global sync event bus singleton
func SyncBus() Bus {
	globalBusLock.Lock()
	defer globalBusLock.Unlock()

	return busInstance
}
