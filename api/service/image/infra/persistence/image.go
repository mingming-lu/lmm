package persistence

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"lmm/api/service/image/domain/model"
	"lmm/api/storage/db"
	"lmm/api/storage/uploader"
)

// ImageStorage s a ImageRepository implementation
type ImageStorage struct {
	db       db.DB
	uploader uploader.Uploader
}

// NewImageStorage creates new ImageStorage
func NewImageStorage(db db.DB, uploader uploader.Uploader) *ImageStorage {
	return &ImageStorage{db: db, uploader: uploader}
}

// Save implementation
func (s *ImageStorage) Save(c context.Context, image *model.Image) error {
	tx, err := s.db.Begin(c, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(c, `insert into image (name, user, created_at) values (?, ?, ?)`)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	if _, err := stmt.Exec(c, image.Name(), image.Uploader().ID(), time.Now()); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	if err := s.uploader.Upload(image.Name(), image.Data()); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	return tx.Commit()
}

// Remove implementation
func (s *ImageStorage) Remove(c context.Context, image *model.Image) error {
	panic("not implemented")
}
