package application

import (
	"context"

	"lmm/api/service/asset/domain/model"
)

type CacheService interface {
	FetchPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, bool)
	StorePhotos(c context.Context, page, perPage uint, photos []*model.PhotoDescriptor) error
	ClearPhotos(c context.Context) error
}

type NopCacheService struct {
	CacheService
}

func (cache *NopCacheService) FetchPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, bool) {
	return nil, false
}
