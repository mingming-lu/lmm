package persistence

import (
	"context"
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/asset/usecase"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

var testAssetRepository usecase.AssetRepository = &AssetDataStore{}

type AssetDataStore struct {
	dataStore *datastore.Client
	transaction.Manager
}

func NewAssetDataStore(dsStore *datastore.Client) *AssetDataStore {
	return &AssetDataStore{
		dataStore: dsStore,
		Manager:   dsUtil.NewTransactionManager(dsStore),
	}
}

type asset struct {
	CreatedAt time.Time `datastore:"CreatedAt"`
	Filename  string    `datastore:"Filename"`
	Type      string    `datastore:"Type"`
}

func (s *AssetDataStore) Save(c context.Context, model *usecase.Asset) error {
	userKey := datastore.IDKey(dsUtil.UserKind, model.UserID, nil)
	assetKey := datastore.IncompleteKey(dsUtil.AssetKind, userKey)

	_, err := dsUtil.MustTransaction(c).Put(assetKey, &asset{
		CreatedAt: model.UploadedAt,
		Filename:  model.Filename,
		Type:      model.Type,
	})
	return err
}

func (s *AssetDataStore) ListPhotos(c context.Context, count int, cursor string) ([]*usecase.Photo, string, error) {
	q := datastore.NewQuery(dsUtil.AssetKind).Project("Filename").Filter("Type =", "Photo").Order("-CreatedAt").Limit(count)
	dsCursor, err := datastore.DecodeCursor(cursor)
	if err == nil {
		q = q.Start(dsCursor)
	}

	photo := asset{}
	photos := make([]*usecase.Photo, 0)

	iter := s.dataStore.Run(c, q)

Iteration:
	for {
		if _, err := iter.Next(&photo); err != nil {
			if err == iterator.Done {
				break Iteration
			}
			return nil, "", errors.Wrap(err, "internal error")
		}

		photos = append(photos, &usecase.Photo{URL: publicURLBase + photo.Filename})
	}

	nextCursor, err := iter.Cursor()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get datastore cursor")
	}

	return photos, nextCursor.String(), nil
}
