package subscription

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type APIAccessLog struct {
	id      string
	storage *datastore.Client
}

func NewAPIAccessLog(dsClient *datastore.Client) *APIAccessLog {
	id := os.Getenv("GCP_PUBSUB_SUBSCRIPTION_API_ACCESS_LOG")
	return &APIAccessLog{
		id:      id,
		storage: dsClient,
	}
}

func (s *APIAccessLog) ID() string {
	return s.id
}

func (s *APIAccessLog) Subscriber() func(context.Context, *pubsub.Message) {
	type accessLog struct {
		ZapBase
		Status       int    `datastore:"status"       json:"status"`
		RequestID    string `datastore:"requestID"    json:"request_id"`
		ClientIP     string `datastore:"clientIP"     json:"client_ip"`
		ForwardedFor string `datastore:"forwardedFor" json:"forwarded_for"`
		UserAgent    string `datastore:"userAgent"    json:"ua"`
		Method       string `datastore:"method"       json:"method"`
		Host         string `datastore:"host"         json:"host"`
		URI          string `datastore:"uri"          json:"uri"`
		Latency      string `datastore:"latency"      json:"latency"`
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
