package subscription

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type APILog struct {
	id      string
	storage *datastore.Client
}

func NewAPILog(dsClient *datastore.Client) *APILog {
	id := os.Getenv("GCP_PUBSUB_SUBSCRIPTION_API_LOG")
	return &APILog{
		id:      id,
		storage: dsClient,
	}
}

func (s *APILog) ID() string {
	return s.id
}

func (s *APILog) Subscriber() func(context.Context, *pubsub.Message) {
	type accessLog struct {
		ZapBase
		RequestID string `datastore:"requestID" json:"request_id"`
	}

	return func(c context.Context, msg *pubsub.Message) {
		al := accessLog{}
		if err := json.Unmarshal(msg.Data, &al); err != nil {
			zap.L().Error(err.Error(),
				zap.String("data", string(msg.Data[:])),
			)
		}

		k := datastore.IncompleteKey(al.LoggerName, nil)
		if key, err := s.storage.Put(c, k, &al); err != nil {
			zap.L().Error(err.Error(),
				zap.String("data", string(msg.Data[:])),
			)
		} else {
			msg.Ack()
			zap.L().Info("saved to datastore",
				zap.String("kind", s.ID()),
				zap.Int64("id", key.ID),
				zap.String("request_id", al.RequestID),
			)
		}
	}
}
