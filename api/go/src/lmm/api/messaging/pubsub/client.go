package pubsub

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func NewClient() (*pubsub.Client, error) {
	return pubsub.NewClient(
		context.Background(),
		os.Getenv("GCP_PROJECT_ID"),
		option.WithCredentialsFile("/gcp/credentials/service_account.json"),
	)
}
