package persistence

import (
	"sort"
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	dsEntity "lmm/api/service/article/port/adapter/persistence/internal/datastore"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

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

func (s *ArticleDataStore) buildArticleKey(articleID, authorID int64) *datastore.Key {
	userKey := datastore.IDKey(dsUtil.UserKind, authorID, nil)

	return datastore.IDKey(dsUtil.ArticleKind, articleID, userKey)
}

func (s *ArticleDataStore) NextID(tx transaction.Transaction, authorID int64) (*model.ArticleID, error) {
	key := datastore.IncompleteKey(dsUtil.ArticleKind, datastore.IDKey(dsUtil.UserKind, authorID, nil))
	keys, err := s.dataStore.AllocateIDs(tx, []*datastore.Key{key})
	if err != nil || len(keys) == 0 {
		return nil, errors.Wrap(err, "failed to allocate new article key")
	}

	return model.NewArticleID(keys[0].Encode()), nil
}

// Save saves article into datastore
func (s *ArticleDataStore) Save(tx transaction.Transaction, model *model.Article) error {
	articleKey := dsUtil.MustKey(model.ID().String())

	dstx := dsUtil.MustTransaction(tx)

	// save article
	if _, err := dstx.Mutate(datastore.NewUpsert(articleKey, &dsEntity.Article{
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
	tags := make([]*dsEntity.Tag, len(model.Content().Tags()), len(model.Content().Tags()))
	for i, model := range model.Content().Tags() {
		tagKeys = append(tagKeys, datastore.IncompleteKey(dsUtil.ArticleTagKind, articleKey))
		tags[i] = &dsEntity.Tag{
			Name:      model.Name(),
			Order:     int(model.Order()),
			CreatedAt: time.Now(),
		}
	}

	// save all tags
	if _, err := dstx.PutMulti(tagKeys, tags); err != nil {
		return errors.Wrap(err, "failed to put tags into datastore")
	}

	return nil
}

func (s *ArticleDataStore) FindByID(tx transaction.Transaction, id *model.ArticleID) (*model.Article, error) {
	articleKey, err := datastore.DecodeKey(id.String())
	if err != nil {
		return nil, errors.Wrapf(domain.ErrNoSuchArticle, "%s: %s", err.Error(), id.String())
	}

	dsTx := dsUtil.MustTransaction(tx)
	data := dsEntity.Article{}
	if err := dsTx.Get(articleKey, &data); err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchArticle, err.Error())
	}

	fetchTagsQuery := datastore.NewQuery(dsUtil.ArticleTagKind).Ancestor(articleKey).Transaction(dsTx)
	var tags []*dsEntity.Tag
	if _, err := s.dataStore.GetAll(tx, fetchTagsQuery, &tags); err != nil {
		return nil, errors.Wrap(err, "failed to get article tags")
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Order < tags[j].Order
	})

	content, err := model.NewContent(data.Title, data.Body, func() []string {
		ss := make([]string, len(tags), len(tags))
		for i, t := range tags {
			ss[i] = t.Name
		}
		return ss
	}())
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	author := model.NewAuthor(articleKey.Parent.ID)

	return model.NewArticle(id, author, content, data.CreatedAt, data.LastModified), nil
}

func (s *ArticleDataStore) Remove(tx transaction.Transaction, id *model.ArticleID) error {
	panic("not implemented")
}

func (s *ArticleDataStore) ViewArticle(tx transaction.Transaction, id string) (*model.Article, error) {
	return s.FindByID(tx, model.NewArticleID(id))
}

func (s *ArticleDataStore) ViewArticles(tx transaction.Transaction, count, page int, filter *model.ArticlesFilter) (*model.ArticleListView, error) {
	if filter != nil && filter.Tag != "" {
		return s.viewArticlesFilteredByTag(tx, count, page, filter.Tag)
	}

	return s.viewAllArticles(tx, count, page)
}

func (s *ArticleDataStore) viewAllArticles(tx transaction.Transaction, count, page int) (*model.ArticleListView, error) {
	counting := datastore.NewQuery(dsUtil.ArticleKind)
	paging := datastore.NewQuery(dsUtil.ArticleKind).Project("__key__", "Title", "CreatedAt").Order("-CreatedAt").Limit(count + 1).Offset((page - 1) * count)

	total, err := s.dataStore.Count(tx, counting)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total number of articles")
	}

	var entities []*dsEntity.ArticleItem
	keys, err := s.dataStore.GetAll(tx, paging, &entities)
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	hasNextPage := false
	if len(entities) > int(count) {
		hasNextPage = true
		entities = entities[:int(count)]
	}

	items := make([]*model.ArticleListViewItem, len(entities), len(entities))
	for i, entity := range entities {
		id := model.NewArticleID(keys[i].Encode())
		item, err := model.NewArticleListViewItem(id, entity.Title, time.Unix(entity.CreatedAt/dsUtil.UnixFactor, 0))
		if err != nil {
			return nil, errors.Wrap(err, "internal error")
		}
		items[i] = item
	}

	return model.NewArticleListView(items, "", page, count, total, hasNextPage), nil
}

func (s *ArticleDataStore) viewArticlesFilteredByTag(tx transaction.Transaction, count, page int, tag string) (*model.ArticleListView, error) {
	dstx := dsUtil.MustTransaction(tx)

	counting := datastore.NewQuery(dsUtil.ArticleTagKind).Filter("Name =", tag)
	paging := datastore.NewQuery(dsUtil.ArticleTagKind).Filter("Name =", tag).KeysOnly().Order("-CreatedAt").Limit(count + 1).Offset((page - 1) * count)

	total, err := s.dataStore.Count(tx, counting)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total number of articles")
	}

	keys, err := s.dataStore.GetAll(tx, paging, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count tags")
	}

	articleKeys := make([]*datastore.Key, len(keys))
	for i := range keys {
		articleKeys[i] = keys[i].Parent
	}

	articles := make([]*dsEntity.Article, len(articleKeys))
	if err := dstx.GetMulti(articleKeys, articles); err != nil {
		return nil, errors.Wrap(err, "failed to get articles")
	}

	items := make([]*model.ArticleListViewItem, len(articles), len(articles))
	for i := range items {
		id := model.NewArticleID(articleKeys[i].Encode())
		m, err := model.NewArticleListViewItem(id, articles[i].Title, articles[i].CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "internal error")
		}
		items[i] = m
	}

	hasNextPage := false
	if len(items) > int(count) {
		hasNextPage = true
		items = items[:int(count)]
	}

	return model.NewArticleListView(items, tag, page, count, total, hasNextPage), nil
}

func (s *ArticleDataStore) ViewAllTags(tx transaction.Transaction) ([]*model.TagView, error) {
	q := datastore.NewQuery(dsUtil.ArticleTagKind).Project("Name").DistinctOn("Name").Order("Name")

	var t dsEntity.Tag
	items := make([]*model.TagView, 0)

	iter := s.dataStore.Run(tx, q)
	for {
		_, err := iter.Next(&t)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "internal error: invalid tag")
		}
		cq := datastore.NewQuery(dsUtil.ArticleTagKind).KeysOnly().Filter("Name =", t.Name)
		c, err := s.dataStore.Count(tx, cq)
		if err != nil {
			return nil, errors.Wrapf(err, "infra: error on counting the number of tag named %s", t.Name)
		}
		items = append(items, model.NewTagView(t.Name, c))
	}

	return items, nil
}
