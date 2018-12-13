package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"lmm/logging/pubsub/subscription"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
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

	opts := []option.ClientOption{
		option.WithCredentialsFile("/gcp/credentials/service_account.json"),
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func main() {
	defer pubsubClient.Close()
	defer dataStoreClient.Close()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	term := newSigC()
	term.Register(os.Interrupt, syscall.SIGTERM)

	subscriptions := []subscription.Subscription{
		subscription.NewAPILog(dataStoreClient),
		subscription.NewAPIAccessLog(dataStoreClient),
	}

	for i := range subscriptions {
		go func(i int, pbc *pubsub.Client) {
			if err := pbc.Subscription(subscriptions[i].ID()).Receive(mainCtx, subscriptions[i].Subscriber()); err != nil {
				logger.Error(err.Error())
				term.Close()
			}
		}(i, pubsubClient)
	}

	time.Sleep(3 * time.Second)
	term.Wait()
	term.Close()
	logger.Info("terminating...")
}

type sigC struct {
	c     chan os.Signal
	close sync.Once
}

func newSigC() *sigC {
	return &sigC{
		c: make(chan os.Signal, 1),
	}
}

func (c *sigC) Register(sig ...os.Signal) {
	signal.Notify(c.c, sig...)
}

func (c *sigC) Wait() os.Signal {
	return <-c.c
}

func (c *sigC) Close() {
	c.close.Do(func() {
		signal.Stop(c.c)
		close(c.c)
	})
}
