package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"lmm/api/messaging"
)

// Client is a rabbimq client
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClient creates a new rabbitmq client
func NewClient() *Client {
	conn, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &Client{
		conn:    conn,
		channel: channel,
	}
}

// Publish published events
func (c *Client) Publish(ctx context.Context, e messaging.Event) error {
	msg, ok := e.(*Event)
	if !ok {
		zap.L().Panic("not a rabbitmq event")
	}
	return c.channel.Publish("", e.Topic(), true, false, msg.Message())
}

// Close closes rabbitmq connection
func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	if err := c.channel.Close(); err != nil {
		return err
	}
	return nil
}
