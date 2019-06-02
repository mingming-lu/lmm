package persistence

import (
	"context"
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/asset/usecase"

	"cloud.google.com/go/datastore"
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

type assert struct {
	CreatedAt time.Time `datastore:"CreatedAt"`
	Filename  string    `datastore:"Filename"`
	Type      string    `datastore:"Type"`
}

func (s *AssetDataStore) Save(c context.Context, model *usecase.Asset) error {
	userKey := datastore.IDKey(dsUtil.UserKind, model.UserID, nil)
	assetKey := datastore.IncompleteKey(dsUtil.AssetKind, userKey)

	_, err := dsUtil.MustTransaction(c).Put(assetKey, &assert{
		CreatedAt: model.UploadedAt,
		Filename:  model.Filename,
		Type:      model.Type,
	})
	return err
}
