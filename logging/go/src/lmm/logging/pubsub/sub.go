package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func main() {
	client, err := pubsub.NewClient(
		context.Background(),
		os.Getenv("GCP_PROJECT_ID"),
		option.WithCredentialsFile("/gcp/credentials/service_account.json"),
	)
	defer client.Close()
	if err != nil {
		panic(err)
	}

	err = client.Subscription(os.Getenv("GCP_PROJECT_SUBSCRIPTION_ID")).
		Receive(context.Background(), func(c context.Context, msg *pubsub.Message) {
			fmt.Print(string(msg.Data[:]))
			msg.Ack()
		})
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("OK")
	select {}
}
