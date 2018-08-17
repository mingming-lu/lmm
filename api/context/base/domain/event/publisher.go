package event

import (
	"errors"
	"sync"
)

var (
	ErrEmptyEventTopic = errors.New("event topic cannot be empty")
)

var (
	mux       sync.RWMutex
	publisher *domainEventPublisher
)

type DomainEventPublisher interface {
	Publish(DomainEvent) error
	Subscribe(DomainEvent, DomainEventSubscriber) error
}

type domainEventPublisher struct {
	subscriberMap map[string][]DomainEventSubscriber
}

func DomainEventPublisherInstance() DomainEventPublisher {
	return publisher
}

func (p *domainEventPublisher) Publish(e DomainEvent) error {
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

func (p *domainEventPublisher) Subscribe(e DomainEvent, sub DomainEventSubscriber) error {
	mux.Lock()
	defer mux.Unlock()

	if e.Topic() == "" {
		return ErrEmptyEventTopic
	}

	subs, ok := p.subscriberMap[e.Topic()]
	if !ok {
		subs = make([]DomainEventSubscriber, 1, 1)
	}
	subs = append(subs, sub)
	p.subscriberMap[e.Topic()] = subs
	return nil
}
