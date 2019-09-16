package internal

import (
	"time"

	"cloud.google.com/go/datastore"
)

type Article struct {
	Title        string    `datastore:"Title"`
	Body         string    `datastore:"Body,noindex"`
	CreatedAt    time.Time `datastore:"CreatedAt"`
	LastModified time.Time `datastore:"LastModified,noindex"`
}

type Tag struct {
	ID        *datastore.Key `datastore:"__key__"`
	Name      string         `datastore:"Name"`
	Order     int            `datastore:"Order"`
	CreatedAt time.Time      `datastore:"CreatedAt"`
}

type ArticleItem struct {
	Title     string `datastore:"Title"`
	CreatedAt int64  `datastore:"CreatedAt"`
}
