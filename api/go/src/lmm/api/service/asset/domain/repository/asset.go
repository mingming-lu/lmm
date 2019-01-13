package repository

import (
	"context"

	"lmm/api/service/asset/domain/model"
)

// AssetRepository provides a interface to deal with persistence of asset
type AssetRepository interface {
	FindByName(c context.Context, name string) (*model.AssetDescriptor, error)
	Save(c context.Context, asset *model.Asset) error
	Remove(c context.Context, asset *model.Asset) error
}
