package pubsub

import (
	"context"
	"io"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"

	"lmm/api/messaging"
)

type PubSubTopicPublisher struct {
	messaging.Publisher
	io.Writer

	topic        *pubsub.Topic
	writeContext func() context.Context
}

func NewPubSubTopicPublisher(topic *pubsub.Topic, writeContext func() context.Context) *PubSubTopicPublisher {
	return &PubSubTopicPublisher{
		topic:        topic,
		writeContext: writeContext,
	}
}

func (p *PubSubTopicPublisher) Publish(ctx context.Context, evt messaging.Event) error {
	e, ok := evt.(*Event)
	if !ok {
		return errors.New("not a pubsub event")
	}

	res := p.topic.Publish(ctx, e.msg)
	<-res.Ready()
	_, err := res.Get(ctx)
	return err
}

func (p *PubSubTopicPublisher) Write(b []byte) (int, error) {
	err := p.Publish(p.writeContext(), NewEvent(p.topic.ID(), b))

	return len(b), err
}
