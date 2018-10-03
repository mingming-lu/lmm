package service

import (
	"context"

	"lmm/api/service/asset/domain/model"
)

// AssetFinder provides interface to find assets
type AssetFinder interface {
	FindAllImages(c context.Context, limit, cursor uint) (*model.ImageCollection, error)
	FindAllPhotos(c context.Context, limit, cursor uint) (*model.PhotoCollection, error)
}
