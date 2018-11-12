package event

import (
	"context"
	"sync"

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

	allErrors := []error{}

	for _, handler := range b.handlers[e.Topic()] {
		if err := handler(c, e); err != nil {
			allErrors = append(allErrors, err)
		}
	}

	if len(allErrors) > 0 {
		err := ErrEventHandleFailed
		for _, e := range allErrors {
			errors.Wrap(err, e.Error())
		}
		return err
	}
	return nil
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

// GlobalBus ger event bus singleton
func GlobalBus() Bus {
	globalBusLock.Lock()
	defer globalBusLock.Unlock()

	return busInstance
}
