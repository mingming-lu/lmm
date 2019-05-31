package persistence

import (
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/viewer"

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
	if err != nil {
		return nil, errors.Wrap(err, "failed to allocate new article key")
	}

	return model.NewArticleID(keys[0].ID, authorID), nil
}

type article struct {
	ID           *datastore.Key `datastore:"__key__"`
	LinkName     string         `datastore:"LinkName"`
	Title        string         `datastore:"Title"`
	Body         string         `datastore:"Body"`
	CreatedAt    time.Time      `datastore:"CreatedAt"`
	LastModified time.Time      `datastore:"LastModified"`
}

type tag struct {
	Name string `datastore:"Name"`
}

// Save saves article into datastore
func (s *ArticleDataStore) Save(tx transaction.Transaction, model *model.Article) error {
	articleKey := s.buildArticleKey(model.ID().ID(), model.ID().AuthorID())

	dstx := dsUtil.MustTransaction(tx)

	// save article
	if _, err := dstx.Mutate(datastore.NewUpsert(articleKey, &article{
		LinkName:     model.LinkName(),
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
	articleKey := s.buildArticleKey(id.ID(), id.AuthorID())

	dsTx := dsUtil.MustTransaction(tx)
	data := article{}
	if err := dsTx.Get(articleKey, &data); err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchArticle, err.Error())
	}

	fetchTagsQuery := datastore.NewQuery(dsUtil.ArticleTagKind).Ancestor(articleKey).Transaction(dsTx)
	var tags []*tag
	if _, err := s.dataStore.GetAll(tx, fetchTagsQuery, &tags); err != nil {
		return nil, errors.Wrap(err, "failed to get article tags")
	}

	text, err := model.NewText(data.Title, data.Body)
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	content := model.NewContent(text, func() []string {
		ss := make([]string, len(tags), len(tags))
		for i, tag := range tags {
			ss[i] = tag.Name
		}
		return ss
	}())

	return model.NewArticle(id, data.LinkName, content, data.CreatedAt, data.LastModified), nil
}

func (s *ArticleDataStore) Remove(tx transaction.Transaction, id *model.ArticleID) error {
	panic("not implemented")
}

func (s *ArticleDataStore) ViewArticle(tx transaction.Transaction, linkName string) (*model.Article, error) {
	q := datastore.NewQuery(dsUtil.ArticleKind).KeysOnly().Filter("LinkName =", linkName).Limit(1)

	keys, err := s.dataStore.GetAll(tx, q, nil)
	if err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchArticle, err.Error())
	}
	if len(keys) == 0 {
		return nil, domain.ErrNoSuchArticle
	}

	k := keys[0]

	return s.FindByID(tx, model.NewArticleID(k.ID, k.Parent.ID))
}

type articleItem struct {
	ID        *datastore.Key `datastore:"__key__"`
	LinkName  string         `datastore:"LinkName"`
	Title     string         `datastore:"Title"`
	CreatedAt int64          `datastore:"CreatedAt"`
}

func (s *ArticleDataStore) ViewArticles(tx transaction.Transaction, count, page int, filter *viewer.ArticlesFilter) (*model.ArticleListView, error) {
	counting := datastore.NewQuery(dsUtil.ArticleKind)
	paging := datastore.NewQuery(dsUtil.ArticleKind).Project("__key__", "LinkName", "Title", "CreatedAt").Order("-CreatedAt").Limit(count + 1).Offset((page - 1) * count)

	total, err := s.dataStore.Count(tx, counting)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total number of articles")
	}

	var entities []*articleItem
	if _, err := s.dataStore.GetAll(tx, paging, &entities); err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	hasNextPage := false
	if len(entities) > int(count) {
		hasNextPage = true
		entities = entities[:int(count)]
	}

	items := make([]*model.ArticleListViewItem, len(entities), len(entities))
	for i, entity := range entities {
		item, err := model.NewArticleListViewItem(entity.ID.ID, entity.LinkName, entity.Title, time.Unix(entity.CreatedAt/dsUtil.UnixFactor, 0))
		if err != nil {
			return nil, errors.Wrap(err, "internal error")
		}
		items[i] = item
	}

	return model.NewArticleListView(items, page, count, total, hasNextPage), nil
}

func (s *ArticleDataStore) ViewAllTags(tx transaction.Transaction) ([]*model.TagView, error) {
	q := datastore.NewQuery(dsUtil.ArticleTagKind).Project("Name").Order("Name").Distinct()

	var t tag
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
		items = append(items, model.NewTagView(t.Name))
	}

	return items, nil
}
