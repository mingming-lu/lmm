package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"lmm/api/messaging"
	"lmm/api/util"
)

// Client is a rabbimq client
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// DefaultClient tries to connect to rabbitmq
func DefaultClient() *Client {
	user := os.Getenv("RABBIT_USER")
	if user == "" {
		user = "guest"
	}

	pass := os.Getenv("RABBIT_PASS")
	if pass == "" {
		pass = "guest"
	}

	host := os.Getenv("RABBIT_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("RABBIT_PORT")
	if port == "" {
		port = "5672"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)

	var (
		client *Client
		err    error
	)

	util.Retry(-1, func() error {
		client, err = NewClient(url)
		if err != nil {
			fmt.Printf(
				"retry connecting to rabbitmq... error: %s, host: %s, port: %s, user: %s.",
				err.Error(), host, port, user,
			)
			<-time.After(5 * time.Second)
		}
		return err
	})
	return client
}

// NewClient creates a new rabbitmq client
func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, channel: channel}, nil
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
