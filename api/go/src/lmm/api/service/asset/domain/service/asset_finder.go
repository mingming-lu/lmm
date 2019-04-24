package service

import (
	"context"

	"lmm/api/service/asset/domain/model"
)

// AssetFinder provides interface to find assets
type AssetFinder interface {
	FindAllImages(context.Context, uint, uint) (*model.ImageCollection, error)
	FindAllPhotos(context.Context, uint, uint) (*model.PhotoCollection, error)
}
