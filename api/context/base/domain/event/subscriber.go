package event

type DomainEventSubscriber interface {
	HandleEvent(e DomainEvent) error
	SubscribedToEventTopic() string
}
