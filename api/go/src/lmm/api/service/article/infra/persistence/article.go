package persistence

import (
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

var testArticleRepo repository.ArticleRepository = &ArticleDataStore{}

type ArticleDataStore struct {
	dataStore *datastore.Client
	transaction.Manager
}

func NewArticleDataStore(dataStore *datastore.Client) *ArticleDataStore {
	return &ArticleDataStore{
		dataStore: dataStore,
		Manager:   dsUtil.NewTransactionManager(dataStore),
	}
}

func (s *ArticleDataStore) NextID(tx transaction.Transaction, authorID int64) (*model.ArticleID, error) {
	key := datastore.IncompleteKey(dsUtil.ArticleKind, datastore.IDKey(dsUtil.UserKind, authorID, nil))
	keys, err := s.dataStore.AllocateIDs(tx, []*datastore.Key{key})
	if err != nil {
		return nil, errors.Wrap(err, "failed to allocate new article key")
	}

	return model.NewArticleID(keys[0].ID), nil
}

type article struct {
	ID           *datastore.Key `datastore:"__key__"`
	Title        string         `datastore:"Title"`
	Body         string         `datastore:"Body"`
	CreatedAt    time.Time      `datastore:"Created_at"`
	LastModified time.Time      `datastore:"LastModified"`
}

type tag struct {
	Name string `datastore:"Name"`
}

// Save saves article into datastore
func (s *ArticleDataStore) Save(tx transaction.Transaction, model *model.Article) error {
	userKey := datastore.IDKey(dsUtil.UserKind, model.Author().ID(), nil)
	articleKey := datastore.IDKey(dsUtil.ArticleKind, int64(*model.ID()), userKey)

	dstx := dsUtil.MustTransaction(tx)

	// save article
	if _, err := dstx.Mutate(datastore.NewUpsert(articleKey, &article{
		Title:        model.Content().Text().Title(),
		Body:         model.Content().Text().Body(),
		CreatedAt:    model.CreatedAt(),
		LastModified: model.LastModified(),
	})); err != nil {
		return errors.Wrap(err, "failed to put article into datastore")
	}

	// get all tag keys by article
	q := datastore.NewQuery(dsUtil.ArticleTagKind).Ancestor(articleKey).KeysOnly().Transaction(dstx)
	tagKeys, err := s.dataStore.GetAll(tx, q, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get article's tags")
	}

	// delete all tags
	if err := dstx.DeleteMulti(tagKeys); err != nil {
		return errors.Wrap(err, "failed to clear article tags")
	}

	tagKeys = tagKeys[:0]
	tags := make([]*tag, len(model.Content().Tags()), len(model.Content().Tags()))
	for i, name := range model.Content().Tags() {
		tagKeys = append(tagKeys, datastore.IncompleteKey(dsUtil.ArticleTagKind, articleKey))
		tags[i] = &tag{Name: name}
	}

	// save all tags
	if _, err := dstx.PutMulti(tagKeys, tags); err != nil {
		return errors.Wrap(err, "failed to put tags into datastore")
	}

	return nil
}

func (s *ArticleDataStore) FindByID(tx transaction.Transaction, id *model.ArticleID) (*model.Article, error) {
	panic("not implemented")
}

func (s *ArticleDataStore) Remove(tx transaction.Transaction, id *model.ArticleID) error {
	panic("not implemented")
}
