package event

type DomainEvent interface {
	Topic() string
}
