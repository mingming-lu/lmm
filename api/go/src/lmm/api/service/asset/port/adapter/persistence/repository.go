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

func (s *AssetDataStore) NextID(c context.Context, userID int64) (*usecase.AssetID, error) {
	parentKey := datastore.IDKey(dsUtil.UserKind, userID, nil)
	pendingKey := datastore.IncompleteKey(dsUtil.AssetKind, parentKey)

	keys, err := s.dataStore.AllocateIDs(c, []*datastore.Key{pendingKey})
	if err != nil {
		return nil, err
	}

	if len(keys) != 1 {
		return nil, errors.New("failed to allocate new asset id")
	}

	assetKey := keys[0]

	return &usecase.AssetID{ID: assetKey.ID, UserID: assetKey.Parent.ID}, nil
}

func (s *AssetDataStore) assetKey(id *usecase.AssetID) *datastore.Key {
	userKey := datastore.IDKey(dsUtil.UserKind, id.UserID, nil)

	return datastore.IDKey(dsUtil.AssetKind, id.ID, userKey)
}

func (s *AssetDataStore) Save(c context.Context, model *usecase.Asset) error {
	_, err := dsUtil.MustTransaction(c).Put(s.assetKey(model.ID), &asset{
		CreatedAt: model.UploadedAt,
		Filename:  model.Filename,
		Type:      model.Type.String(),
	})
	return err
}

type photoTag struct {
	Name  string `datastore:"Name"`
	Order int    `datastore:"Order"`
}

func (s *AssetDataStore) SetPhotoTags(c context.Context, id *usecase.AssetID, tags []string) error {
	assetKey := s.assetKey(id)

	tx := dsUtil.MustTransaction(c)
	q := datastore.NewQuery(dsUtil.PhotoTagKind).Ancestor(assetKey).KeysOnly().Transaction(tx)

	keys, err := s.dataStore.GetAll(c, q, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get photo tag keys")
	}

	if err := tx.DeleteMulti(keys); err != nil {
		return errors.Wrap(err, "failed to delete clear photo tags")
	}

	keys = keys[:0]
	newTags := make([]*photoTag, len(tags), len(tags))

	for i, name := range tags {
		keys = append(keys, datastore.IncompleteKey(dsUtil.PhotoTagKind, assetKey))
		newTags[i] = &photoTag{Name: name, Order: i + 1}
	}

	if _, err := tx.PutMulti(keys, newTags); err != nil {
		return errors.Wrap(err, "failed to save new photo tags")
	}

	return nil
}

func (s *AssetDataStore) Find(c context.Context, id *usecase.AssetID) (*usecase.Asset, error) {
	var model asset
	err := dsUtil.MustTransaction(c).Get(s.assetKey(id), &model)

	return &usecase.Asset{
		Filename:   model.Filename,
		Type:       usecase.AssetTypeFromString(model.Type),
		UploadedAt: model.CreatedAt,
	}, err
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
