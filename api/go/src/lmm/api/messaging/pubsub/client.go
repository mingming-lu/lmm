package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// NewClient creates a new gcp pubsub client
func NewClient(projectID, credentialJSONPath string) (*pubsub.Client, error) {
	return pubsub.NewClient(
		context.Background(),
		projectID,
		option.WithCredentialsFile(credentialJSONPath),
	)
}
