package persistence

import (
	"context"
	"fmt"
	"sort"
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/asset/usecase"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

var _ usecase.AssetRepository = &AssetDataStore{}

type AssetDataStore struct {
	dataStore *datastore.Client
	transaction.Manager
	gcsBucketName string
}

func NewAssetDataStore(
	c context.Context,
	dsStore *datastore.Client,
	gcsBucket *storage.BucketHandle,
) (*AssetDataStore, error) {
	attr, err := gcsBucket.Attrs(c)
	if err != nil {
		return nil, err
	}

	return &AssetDataStore{
		gcsBucketName: attr.Name,
		dataStore:     dsStore,
		Manager:       dsUtil.NewTransactionManager(dsStore),
	}, nil
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

	return usecase.NewAssetID(keys[0].Encode()), nil
}

func (s *AssetDataStore) assetKey(id *usecase.AssetID) (*datastore.Key, error) {
	key, err := datastore.DecodeKey(id.String())
	if err != nil {
		return nil, errors.Wrapf(err, "invalid encoded datastore key: %s", id.String())
	}
	return key, nil
}

func (s *AssetDataStore) Save(c context.Context, model *usecase.Asset) error {
	key, err := s.assetKey(model.ID)
	if err != nil {
		return errors.Wrap(err, "error occurred on save asset")
	}

	if _, err := dsUtil.MustTransaction(c).Put(key, &asset{
		CreatedAt: model.UploadedAt,
		Filename:  model.Filename,
		Type:      model.Type.String(),
	}); err != nil {
		return err
	}
	return nil
}

type photoTag struct {
	Name  string `datastore:"Name"`
	Order int    `datastore:"Order"`
}

func (s *AssetDataStore) SetPhotoTags(c context.Context, id *usecase.AssetID, tags []string) error {
	assetKey, err := s.assetKey(id)
	if err != nil {
		return errors.Wrap(err, "error occurred on set photo tags")
	}

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
	key, err := datastore.DecodeKey(id.String())
	if err != nil {
		return nil, errors.Wrapf(err, "invalid encoded datastore key: %s", id.String())
	}

	var model asset
	if err := dsUtil.MustTransaction(c).Get(key, &model); err != nil {
		return nil, errors.Wrap(err, "internel error: failed to get asset by key")
	}

	return &usecase.Asset{
		ID:         usecase.NewAssetID(key.Encode()),
		UserID:     key.Parent.ID,
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
		key, err := iter.Next(&photo)
		if err != nil {
			if err == iterator.Done {
				break Iteration
			}
			return nil, "", errors.Wrap(err, "internal error")
		}

		tags, err := s.getTagsByAssetKey(c, key)
		if err != nil {
			return nil, "", errors.Wrap(err, "error occurred on getting photo list")
		}

		photos = append(photos, &usecase.Photo{
			ID:   key.Encode(),
			URL:  s.GetPublicURL(c, photo.Filename),
			Tags: tags,
		})
	}

	nextCursor, err := iter.Cursor()
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get datastore cursor")
	}

	return photos, nextCursor.String(), nil
}

func (s *AssetDataStore) GetPublicURL(c context.Context, filename string) string {
	return fmt.Sprintf(templatePublicURL, s.gcsBucketName, filename)
}

func (s *AssetDataStore) GetTagsByPhotoID(c context.Context, id *usecase.AssetID) ([]string, error) {
	key, err := s.assetKey(id)
	if err != nil {
		return nil, errors.Wrap(err, "error occurred on get photo tags")
	}
	return s.getTagsByAssetKey(c, key)
}

func (s *AssetDataStore) getTagsByAssetKey(c context.Context, key *datastore.Key) ([]string, error) {
	var tags []*photoTag
	qt := datastore.NewQuery(dsUtil.PhotoTagKind).Ancestor(key).Transaction(dsUtil.MustTransaction(c))
	if _, err := s.dataStore.GetAll(c, qt, &tags); err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Order < tags[j].Order
	})

	tagNames := func() []string {
		names := make([]string, len(tags), len(tags))
		for i, tag := range tags {
			names[i] = tag.Name
		}
		return names
	}()

	return tagNames, nil
}
