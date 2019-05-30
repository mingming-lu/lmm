package persistence

import (
	"context"
	"testing"

	"lmm/api/clock"
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
			assert.NotZero(t, int64(articleID))

			return nil
		}, nil)
	})

	t.Run("Save", func(t *testing.T) {
		var article *model.Article

		t.Run("Insert", func(t *testing.T) {
			articleDataStore.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				articleID, err := articleDataStore.NextID(tx, 1)
				assert.NoError(t, err)
				assert.NotZero(t, int64(articleID))

				text, err := model.NewText(uuidutil.NewUUID(), uuidutil.NewUUID())
				if err != nil {
					t.Fatal(errors.Wrap(err, "internal error"))
				}

				now := clock.Now()
				article = model.NewArticle(articleID, model.NewAuthor(1), model.NewContent(text, []string{}), now, now)
				if !assert.NoError(t, articleDataStore.Save(tx, article)) || !assert.NotNil(t, article) {
					t.Fatal("failed to save article")
				}

				return nil
			}, nil)
		})

		t.Run("Update", func(t *testing.T) {
			assert.NoError(t, articleDataStore.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				text, err := model.NewText(uuidutil.NewUUID(), uuidutil.NewUUID())
				if err != nil {
					t.Fatal(errors.Wrap(err, "internal error"))
				}
				article.EditContent(model.NewContent(text, []string{"tag1", "tag2"}))

				return articleDataStore.Save(tx, article)
			}, nil))
		})
	})
}
