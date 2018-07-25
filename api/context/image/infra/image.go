package infra

import (
	"lmm/api/context/image/domain"
	"lmm/api/context/image/domain/model"
	"lmm/api/storage"
	"time"
)

type ImageStorage struct {
	db *storage.DB
}

func NewImageStorage(db *storage.DB) *ImageStorage {
	return &ImageStorage{db: db}
}

func (s *ImageStorage) Add(image *model.Image) error {
	stmt := s.db.MustPrepare(`INSERT INTO image (uid, user, created_at) VALUES (?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(image.ID(), image.UserID(), image.CreatedAt())

	if key, _, ok := storage.CheckErrorDuplicate(err); ok {
		switch key {
		case "uid":
			return domain.ErrDuplicateImageID
		}
	}

	return err
}

func (s *ImageStorage) Remove(image *model.Image) error {
	stmt := s.db.MustPrepare(`DELETE FROM image WHERE uid = ?`)
	defer stmt.Close()

	res, err := stmt.Exec(image.ID())

	if rowsAffected, e := res.RowsAffected(); e != nil {
		return e
	} else if rowsAffected == 0 {
		return domain.ErrNoSuchImage
	}

	return err
}

func (s *ImageStorage) FindByID(id string) (*model.Image, error) {
	stmt := s.db.MustPrepare(`SELECT uid, user, created_at FROM image WHERE uid = ?`)
	defer stmt.Close()

	var (
		imageID           string
		imageUploaderID   uint64
		imageUploaderTime time.Time
	)
	if err := stmt.QueryRow(id).Scan(&imageID, &imageUploaderID, &imageUploaderTime); err != nil {
		switch err {
		case storage.ErrNoRows:
			return nil, domain.ErrNoSuchImage
		}
	}

	return model.NewImage(imageID, imageUploaderID, imageUploaderTime), nil
}
