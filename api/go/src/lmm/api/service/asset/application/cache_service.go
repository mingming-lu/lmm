package application

import (
	"context"

	"lmm/api/service/asset/domain/model"
)

// CacheService defines the interface to deal with asset cache
type CacheService interface {
	FetchPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, bool)
	StorePhotos(c context.Context, page, perPage uint, photos []*model.PhotoDescriptor) error
	ClearPhotos(c context.Context) error
}

type nopCacheService struct {
	CacheService
}

var nop = &nopCacheService{}

func (cache *nopCacheService) FetchPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, bool) {
	return nil, false
}

func (cache *nopCacheService) StorePhotos(c context.Context, page, perPage uint, photos []*model.PhotoDescriptor) error {
	return nil
}

func (cache *nopCacheService) ClearPhotos(c context.Context) error {
	return nil
}

// NopCacheService returns nop CacheService implementation
func NopCacheService() CacheService {
	return nop
}
