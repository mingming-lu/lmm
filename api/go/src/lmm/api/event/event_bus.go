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
	defer globalBusLock.RUnlock()

	if e.Topic() == "" {
		return ErrEmptyEventTopic
	}

	cancelCtx, cancel := context.WithCancel(c)
	group, ctx := errgroup.WithContext(cancelCtx)

	for _, handler := range b.handlers[e.Topic()] {
		group.Go(func() error {
			return handler(c, e)
		})
	}

	go func() {
		group.Wait()
		cancel()
	}()

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return group.Wait()
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
