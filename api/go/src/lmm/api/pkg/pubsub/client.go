package pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"lmm/api/messaging"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// Client wraps pubsub client and implements Publisher and Subscriber
type Client struct {
	rwMutex      sync.RWMutex
	pubsubClient *pubsub.Client
	topics       map[string]*pubsub.Topic
}

// NewClient create a new pubsub client
func NewClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	c, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		pubsubClient: c,
		topics:       make(map[string]*pubsub.Topic),
	}, nil
}

// Close closes c
func (c *Client) Close() error {
	c.rwMutex.Lock()
	for _, topic := range c.topics {
		topic.Stop()
	}
	c.rwMutex.Unlock()

	return c.pubsubClient.Close()
}

func (c *Client) getOrCreateTopic(ctx context.Context, name string) (*pubsub.Topic, error) {
	c.rwMutex.RLock()
	if topic, ok := c.topics[name]; ok {
		c.rwMutex.RUnlock()
		return topic, nil
	}

	c.rwMutex.RUnlock()
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	topic := c.pubsubClient.Topic(name)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		c.topics[name] = topic
		return topic, nil
	}

	topic, err = c.pubsubClient.CreateTopic(ctx, name)
	if err != nil {
		return nil, err
	}

	c.topics[name] = topic
	return topic, nil
}

// Publish publishes evt to pub/sub
func (c *Client) Publish(ctx context.Context, evt messaging.Event) error {
	buf := getBuf()

	if err := json.NewEncoder(buf).Encode(evt.Message()); err != nil {
		return fmt.Errorf("failed to encode event message into json: %v", err)
	}

	topic, err := c.getOrCreateTopic(ctx, evt.Topic())
	if err != nil {
		return err
	}

	msg, err := EventToPubSubMessage(evt)
	if err != nil {
		return err
	}

	result := topic.Publish(ctx, msg)

	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("failed on pubsub topic publish: %v", err)
	}

	freeBuf(buf)

	return nil
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func getBuf() *bytes.Buffer {
	obj, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		panic(fmt.Sprintf("expect *bytes.Buffer but got %T", obj))
	}

	obj.Reset()
	return obj
}

func freeBuf(buf *bytes.Buffer) {
	bufPool.Put(buf)
}

type domainEventAdapter struct {
	T string      `json:"t"`
	P time.Time   `json:"p"`
	M interface{} `json:"m"`
}

func (e *domainEventAdapter) Topic() string {
	return e.T
}

func (e *domainEventAdapter) PublishedAt() time.Time {
	return e.P
}

func (e *domainEventAdapter) Message() interface{} {
	return e.M
}

// Subscribe subscribes topic on pub/sub
func (c *Client) Subscribe(ctx context.Context, topic string, handler messaging.EventHandler) (err error) {
	sub := c.pubsubClient.Subscription(topic)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		pubsubTopic, err := c.getOrCreateTopic(ctx, topic)
		if err != nil {
			return err
		}
		sub, err = c.pubsubClient.CreateSubscription(ctx, topic, pubsub.SubscriptionConfig{
			Topic: pubsubTopic,
		})
		if err != nil {
			return err
		}
	}

	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		evt, err := EventFromPubSubMessage(msg)
		if err != nil {
			log.Printf("invalid pubsub message. Error: %s Data: %s",
				err, string(msg.Data[:]),
			)
			return
		}
		err = handler(ctx, evt)
		if err != nil {
			log.Printf("failed to handle pubsub event. Error: %s Data: %s",
				err, string(msg.Data[:]),
			)
			return
		}
		msg.Ack()
	})
}

// EventFromPubSubMessage adapts msg to Event
func EventFromPubSubMessage(msg *pubsub.Message) (messaging.Event, error) {
	buf := getBuf()
	if _, err := buf.Write(msg.Data); err != nil {
		return nil, err
	}

	var obj domainEventAdapter

	if err := json.NewDecoder(bytes.NewBuffer(msg.Data)).Decode(&obj); err != nil {
		return nil, err
	}

	freeBuf(buf)

	return &obj, nil
}

// EventToPubSubMessage adapts evt to Message
func EventToPubSubMessage(evt messaging.Event) (*pubsub.Message, error) {
	obj := domainEventAdapter{
		T: evt.Topic(),
		P: evt.PublishedAt(),
		M: evt.Message(),
	}

	buf := getBuf()

	if err := json.NewEncoder(buf).Encode(obj); err != nil {
		return nil, err
	}

	msg := &pubsub.Message{
		Data: buf.Bytes(),
	}

	freeBuf(buf)

	return msg, nil
}

// ScanEvent scans evt into obj
func ScanEvent(evt messaging.Event, obj interface{}) error {
	b, err := json.Marshal(evt.Message())
	if err != nil {
		return err
	}

	buf := getBuf()
	if _, err := buf.Write(b); err != nil {
		return err
	}

	if err := json.NewDecoder(buf).Decode(obj); err != nil {
		return err
	}

	freeBuf(buf)

	return nil
}
