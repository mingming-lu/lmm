package persistence

import (
	"context"
	"testing"

	"lmm/api/pkg/transaction"
	"lmm/api/service/article/domain/model"
	"lmm/api/util/uuidutil"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestArticleDataStore(t *testing.T) {
	ctx := context.Background()

	dataStore, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to datastore"))
	}

	articleDataStore := NewArticleDataStore(dataStore)

	t.Run("NextID", func(t *testing.T) {
		articleDataStore.RunInTransaction(ctx, func(tx transaction.Transaction) error {
			articleID, err := articleDataStore.NextID(tx, 1)
			assert.NoError(t, err)
			assert.NotZero(t, int64(*articleID))

			return nil
		}, nil)
	})

	t.Run("Save", func(t *testing.T) {
		articleDataStore.RunInTransaction(ctx, func(tx transaction.Transaction) error {
			articleID, err := articleDataStore.NextID(tx, 1)
			assert.NoError(t, err)
			assert.NotZero(t, int64(*articleID))

			text, err := model.NewText(uuidutil.NewUUID(), uuidutil.NewUUID())
			if err != nil {
				t.Fatal(errors.Wrap(err, "internal error"))
			}

			article := model.NewArticle(articleID, model.NewAuthor(1), model.NewContent(text, []string{}), nil)
			assert.NoError(t, articleDataStore.Save(tx, article))

			return nil
		}, nil)
	})
}
