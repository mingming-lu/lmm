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

func (f *assetFetcher) FindAllImages(c context.Context, limit, nextCursor uint) (*model.ImageCollection, error) {
	var query string
	var args []interface{}
	if nextCursor == 0 {
		query = `select id, name from asset where type = 0 order by created_at desc limit ?`
		args = []interface{}{limit + 1}
	} else {
		query = `select id, name from asset where id <= ? and type = 0 order by created_at desc limit ?`
		args = []interface{}{nextCursor, limit + 1}
	}
	stmt := f.db.Prepare(c, query)
	defer stmt.Close()

	rows, err := stmt.Query(c, args...)
	if err != nil {
		if err == db.ErrNoRows {
			return model.NewImageCollection(nil, nil), nil
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

	var nextID *uint
	if uint(len(images)) > limit {
		leading, trailing := images[:limit], images[limit:]
		id := trailing[0].ID()
		nextID = &id
		images = leading
	}

	return model.NewImageCollection(images, nextID), nil
}

func (f *assetFetcher) FindAllPhotos(c context.Context, limit, nextCursor uint) (*model.PhotoCollection, error) {
	var query string
	var args []interface{}
	if nextCursor == 0 {
		query = `select id, name from asset where type = 1 order by created_at desc limit ?`
		args = []interface{}{limit + 1}
	} else {
		query = `select id, name from asset where id <= ? and type = 1 order by created_at desc limit ?`
		args = []interface{}{nextCursor, limit + 1}
	}
	stmt := f.db.Prepare(c, query)
	defer stmt.Close()

	rows, err := stmt.Query(c, args...)
	if err != nil {
		if err == db.ErrNoRows {
			return model.NewPhotoCollection(nil, nil), nil
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

	var nextID *uint
	if uint(len(photos)) > limit {
		leading, trailing := photos[:limit], photos[limit:]
		id := trailing[0].ID()
		nextID = &id
		photos = leading
	}

	return model.NewPhotoCollection(photos, nextID), nil
}
