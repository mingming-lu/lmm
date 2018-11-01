package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()

	if err != nil {
		log.Fatal(err.Error())
	}

	zap.ReplaceGlobals(logger)
}

func connect() *amqp.Connection {
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
		host = "localhost"
	}

	port := os.Getenv("RABBIT_PORT")
	if port == "" {
		port = "5672"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)

	var (
		conn *amqp.Connection
		err  error
	)

	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		zap.L().Warn("retry connecting to rabbitmq...",
			zap.String("error", err.Error()),
			zap.String("host", host),
			zap.String("port", port),
			zap.String("user", user),
		)
		<-time.After(5 * time.Second)
	}

	zap.L().Info("rabbitmq connected",
		zap.String("host", host),
		zap.String("port", port),
		zap.String("user", user),
	)
	return conn
}

// Subscriber subscribes topics
type Subscriber struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewSubscriber (re)tries to connect to rabbitmq until success
func NewSubscriber() (*Subscriber, error) {
	conn := connect()
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Subscriber{conn: conn, channel: ch}, nil
}

// Close closes subscriber's connection
func (s *Subscriber) Close() error {
	if err := s.channel.Close(); err != nil {
		return err
	}

	if err := s.conn.Close(); err != nil {
		return err
	}

	return nil
}

// Event is an alias of amqp.Delivery
type Event = amqp.Delivery

// Subscribe subscribes topics
func (s *Subscriber) Subscribe(topic string, handler func(event Event) error) {
	q, err := s.channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	msgs, err := s.channel.Consume(q.Name, "", false, false, false, false, nil)

	for msg := range msgs {
		err := handler(msg)
		if err != nil {
			zap.L().Error(err.Error())
		} else {
			if err := msg.Ack(false); err != nil {
				zap.L().Error(err.Error())
			}
		}
	}
}

func main() {
	subscriber, err := NewSubscriber()
	if err != nil {
		zap.L().Error(err.Error())
	}
	defer subscriber.Close()

	term := make(chan bool)

	go subscriber.Subscribe("asset.photo.uploaded", func(event Event) error {
		return uploadPhoto(event.MessageId, event.Body)
	})

	go subscriber.Subscribe("asset.image.uploaded", func(event Event) error {
		return uploadImage(event.MessageId, event.Body)
	})

	<-term
}

func uploadImage(name string, data []byte) error {
	return upload("images", name, data)
}

func uploadPhoto(name string, data []byte) error {
	return upload("photos", name, data)
}

func upload(assetType, name string, data []byte) error {
	file, err := os.OpenFile("/static/"+assetType+"/"+name, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return w.Flush()
}
