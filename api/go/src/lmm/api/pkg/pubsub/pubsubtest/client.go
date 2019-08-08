package pubsubtest

import (
	"context"

	"lmm/api/messaging"
	"lmm/api/pkg/pubsub"

	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type TestPubSubClient struct {
	grpcConn     *grpc.ClientConn
	fakeServer   *pstest.Server
	pubsubClient *pubsub.Client
}

// NewClient creates a new TestPubSubClient
func NewClient() *TestPubSubClient {
	fakeServer := pstest.NewServer()

	grpcConn, err := grpc.Dial(fakeServer.Addr, grpc.WithInsecure())
	if err != nil {
		fakeServer.Close()
		panic(err)
	}

	client, err := pubsub.NewClient(context.Background(), "", option.WithGRPCConn(grpcConn))
	if err != nil {
		fakeServer.Close()
		grpcConn.Close()
		panic(err)
	}

	return &TestPubSubClient{
		grpcConn:     grpcConn,
		fakeServer:   fakeServer,
		pubsubClient: client,
	}
}

func (c *TestPubSubClient) Close() {
	c.pubsubClient.Close()
	c.fakeServer.Close()
	c.grpcConn.Close()
}

func (c *TestPubSubClient) Publish(ctx context.Context, evt messaging.Event) error {
	return c.pubsubClient.Publish(ctx, evt)
}

func (c *TestPubSubClient) Subscribe(ctx context.Context, topic string, handler messaging.EventHandler) error {
	return c.pubsubClient.Subscribe(ctx, topic, handler)
}
