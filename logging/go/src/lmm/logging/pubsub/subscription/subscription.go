package subscription

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// Subscription is an abstraction of GCP pubsub subscription
type Subscription interface {
	ID() string
	Subscriber() func(c context.Context, msg *pubsub.Message)
}
