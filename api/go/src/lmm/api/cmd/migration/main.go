package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/datastore"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/service/article/port/adapter/persistence"
)

func main() {
	baseCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dsClient, err := datastore.NewClient(baseCtx, os.Getenv("DATASTORE_PROJECT_ID"))
	if err != nil {
		panic(err)
	}

	defer dsClient.Close()

	articles := make([]*persistence.Article, 0)
	q := datastore.NewQuery(dsUtil.ArticleKind)
	if _, err := dsClient.GetAll(baseCtx, q, &articles); err != nil {
		panic(err)
	}

	mutations := make([]*datastore.Mutation, 0)
	for _, article := range articles {
		if article.PublishedAt.IsZero() {
			article.PublishedAt = article.CreatedAt
			mutations = append(mutations, datastore.NewUpdate(article.Key, article))
		}
	}

	if _, err := dsClient.Mutate(baseCtx, mutations...); err != nil {
		panic(err)
	}
}
