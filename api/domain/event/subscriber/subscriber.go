package subscriber

import "lmm/api/domain/event"

type DomainEventSubscriber interface {
	HandleEvent(e event.DomainEvent) error
}
