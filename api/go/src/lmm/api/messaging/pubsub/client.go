package pubsub

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Client struct {
	*pubsub.Client
}

func NewClient() *Client {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx,
		os.Getenv("GCP_PROJECT_ID"),
		option.WithCredentialsFile("/gcp/credentials/service_account.json"),
	)
	if err != nil {
		panic(err)
	}

	return &Client{
		Client: c,
	}
}

func (c *Client) Publish(ctx context.Context) error {
	res := c.Topic(os.Getenv("GCP_TOPIC_ID")).Publish(ctx, &pubsub.Message{Data: []byte("hey bro")})
	res.Ready()
	_, err := res.Get(ctx)
	return err
}
