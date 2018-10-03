package service

import (
	"context"

	"lmm/api/service/asset/domain/model"
	"lmm/api/service/asset/domain/service"
	"lmm/api/storage/db"
)

// AssetFetcher implements AssetFinder interface
type assetFetcher struct {
	db db.DB
}

// NewAssetFetcher returns an implementation of AssetFinder
func NewAssetFetcher(db db.DB) service.AssetFinder {
	return &assetFetcher{db: db}
}

func (f *assetFetcher) FindAllImages(c context.Context, page, perPage uint) (*model.ImageCollection, error) {
	stmt := f.db.Prepare(c, `select id, name from asset where type = 0 order by created_at desc limit ? offset ?`)
	defer stmt.Close()

	rows, err := stmt.Query(c, perPage+1, (page-1)*perPage)
	if err != nil {
		if err == db.ErrNoRows {
			return model.NewImageCollection(nil, false), nil
		}
		return nil, err
	}
	defer rows.Close()

	var (
		id   uint
		name string
	)
	images := make([]*model.ImageDescriptor, 0)
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		images = append(images, model.NewImageDescriptor(id, name))
	}

	hasNextPage := false
	if uint(len(images)) > perPage {
		images = images[:perPage]
		hasNextPage = true
	}

	return model.NewImageCollection(images, hasNextPage), nil
}

func (f *assetFetcher) FindAllPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, error) {
	stmt := f.db.Prepare(c, `select id, name from asset where type = 1 order by created_at desc limit ? offset ?`)
	defer stmt.Close()

	rows, err := stmt.Query(c, perPage+1, (page-1)*perPage)
	if err != nil {
		if err == db.ErrNoRows {
			return model.NewPhotoCollection(nil, false), nil
		}
		return nil, err
	}
	defer rows.Close()

	var (
		id   uint
		name string
	)
	photos := make([]*model.PhotoDescriptor, 0)
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		photos = append(photos, model.NewPhotoDescriptor(id, name))
	}

	hasNextPage := false
	if uint(len(photos)) > perPage {
		photos = photos[:perPage]
		hasNextPage = true
	}

	return model.NewPhotoCollection(photos, hasNextPage), nil
}
