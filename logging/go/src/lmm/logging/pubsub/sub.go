package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
)

var (
	projectID            string
	dataStoreLoggingKind string

	dataStoreClient *datastore.Client
	pubsubClient    *pubsub.Client

	logger *zap.Logger
)

func init() {
	var err error

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}

	projectID = os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		logger.Panic("empty project id")
	}
	logger.Info("gcp project id found", zap.String("project_id", projectID))

	dataStoreLoggingKind = os.Getenv("GCP_DATASTORE_LOGGING_KIND")
	if dataStoreLoggingKind == "" {
		logger.Panic("empty kind")
	}
	logger.Info("gcp datastore kind found", zap.String("datastore_kind", dataStoreLoggingKind))

	opts := []option.ClientOption{
		option.WithCredentialsFile("/gcp/credentials/service_account.json"),
	}

	c, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	g, c := errgroup.WithContext(c)

	g.Go(func() error {
		client, err := pubsub.NewClient(c, projectID, opts...)
		if err != nil {
			return err
		}
		pubsubClient = client
		logger.Info("connected to gcp pub/sub")
		return nil
	})

	g.Go(func() error {
		client, err := datastore.NewClient(c, projectID, opts...)
		if err != nil {
			return err
		}
		dataStoreClient = client
		logger.Info("connected to gcp datastore")
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Panic(err.Error())
	}
}

type accessLog struct {
	Level        string `datastore:"level"        json:"level"`
	TimeStamp    string `datastore:"time"         json:"ts"`
	LoggerName   string `datastore:"name"         json:"logger"`
	Message      string `datastore:"message"      json:"msg"`
	Status       int    `datastore:"status"       json:"status"`
	RequestID    string `datastore:"request_id"   json:"request_id"`
	ClientIP     string `datastore:"clientIP"     json:"client_ip"`
	ForwardedFor string `datastore:"forwardedFor" json:"forwarded_for"`
	UserAgent    string `datastore:"userAgent"    json:"ua"`
	Method       string `datastore:"method"       json:"method"`
	Host         string `datastore:"host"         json:"host"`
	URI          string `datastore:"uri"          json:"uri"`
	Latency      string `datastore:"latency"      json:"latency"`
}

func main() {
	defer pubsubClient.Close()
	defer dataStoreClient.Close()

	go func() {
		loggingSubID := os.Getenv("GCP_PUBSUB_LOGGING_SUBSCRIPTION_ID")
		if loggingSubID == "" {
			logger.Panic("empty subscription id")
		}
		logger.Info("listen to pub/sub subscription", zap.String("subscription_id", loggingSubID))

		err := pubsubClient.Subscription(loggingSubID).
			Receive(context.Background(), func(c context.Context, msg *pubsub.Message) {
				al := accessLog{}
				if err := json.Unmarshal(msg.Data, &al); err != nil {
					logger.Error(err.Error(),
						zap.String("data", string(msg.Data[:])),
					)
				}

				k := datastore.IncompleteKey(dataStoreLoggingKind, nil)
				if key, err := dataStoreClient.Put(c, k, &al); err != nil {
					logger.Error(err.Error(),
						zap.String("data", string(msg.Data[:])),
					)
				} else {
					msg.Ack()
					logger.Info("saved to datastore",
						zap.Int64("id", key.ID),
						zap.String("request_id", al.RequestID),
					)
				}
			})
		if err != nil {
			logger.Panic(err.Error())
		}
	}()

	select{}
}
