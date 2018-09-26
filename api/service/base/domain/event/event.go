package event

import "time"

type DomainEvent interface {
	Topic() string
	OccurredAt() time.Time
}
