package subscriber

import "lmm/api/context/base/domain/event"

type DomainEventSubscriber interface {
	HandleEvent(e event.DomainEvent) error
}
