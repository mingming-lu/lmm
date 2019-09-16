package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

type Tag struct {
	ID        *datastore.Key `datastore:"__key__"`
	Name      string         `datastore:"Name"`
	Order     int            `datastore:"Order"`
	CreatedAt time.Time      `datastore:"CreatedAt"`
}

type Article struct {
	ID           *datastore.Key `datastore:"__key__"`
	Title        string         `datastore:"Title"`
	Body         string         `datastore:"Body,noindex"`
	CreatedAt    time.Time      `datastore:"CreatedAt"`
	LastModified time.Time      `datastore:"LastModified,noindex"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dsClient, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if err != nil {
		panic(err)
	}

	cmt, err := dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		q := datastore.NewQuery("ArticleTag")
		var tags []*Tag
		if _, err := dsClient.GetAll(ctx, q, &tags); err != nil {
			return err
		}
		for i := range tags {
			if tags[i].CreatedAt.IsZero() {
				var article Article
				if err := tx.Get(tags[i].ID.Parent, &article); err != nil {
					return err
				}
				tags[i].CreatedAt = article.CreatedAt
				if _, err := tx.Mutate(
					datastore.NewUpdate(tags[i].ID, tags[i]),
				); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	log.Printf("committed: %#v", cmt)
}
