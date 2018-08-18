package publisher

import (
	"errors"
	"sync"

	"lmm/api/context/base/domain/event"
	"lmm/api/context/base/domain/event/subscriber"
)

var (
	ErrEmptyEventTopic = errors.New("event topic cannot be empty")
)

var (
	mux       sync.RWMutex
	publisher *domainEventPublisher
)

type DomainEventPublisher interface {
	Publish(event.DomainEvent) error
	Subscribe(event.DomainEvent, subscriber.DomainEventSubscriber) error
}

type domainEventPublisher struct {
	subscriberMap map[string][]subscriber.DomainEventSubscriber
}

func Instance() DomainEventPublisher {
	return publisher
}

func (p *domainEventPublisher) Publish(e event.DomainEvent) error {
	subscribers, ok := p.subscriberMap[e.Topic()]
	if !ok {
		return errors.New("not subscribed topic: '" + e.Topic() + "'")
	}
	for _, subscriber := range subscribers {
		if err := subscriber.HandleEvent(e); err != nil {
			return err
		}
	}
	return nil
}

func (p *domainEventPublisher) Subscribe(e event.DomainEvent, sub subscriber.DomainEventSubscriber) error {
	mux.Lock()
	defer mux.Unlock()

	if e.Topic() == "" {
		return ErrEmptyEventTopic
	}

	subs, ok := p.subscriberMap[e.Topic()]
	if !ok {
		subs = make([]subscriber.DomainEventSubscriber, 1, 1)
	}
	subs = append(subs, sub)
	p.subscriberMap[e.Topic()] = subs
	return nil
}
