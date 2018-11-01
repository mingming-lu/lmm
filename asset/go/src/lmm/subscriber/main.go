package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
		log.Println(err.Error())
		<-time.After(5 * time.Second)
	}

	return conn
}

func main() {
	conn := connect()
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"asset.photo.uploaded", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.MessageId)
			err := upload(msg.MessageId, msg.Body)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func upload(name string, data []byte) error {
	file, err := os.OpenFile("/static/photos/"+name, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
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
