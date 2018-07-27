package infra

import (
	"lmm/api/context/image/domain"
	"lmm/api/context/image/domain/model"
	"lmm/api/storage"
	"time"
)

type ImageStorage struct {
	db               *storage.DB
	staticRepository storage.StaticRepository
}

func NewImageStorage(db *storage.DB) *ImageStorage {
	return &ImageStorage{db: db}
}

func (s *ImageStorage) SetStaticRepository(repo storage.StaticRepository) {
	s.staticRepository = repo
}

func (s *ImageStorage) Add(image *model.ImageWithData) error {
	stmt := s.db.MustPrepare(`INSERT INTO image (uid, user, created_at) VALUES (?, ?, ?)`)
	defer stmt.Close()

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(image.ID(), image.UserID(), image.CreatedAt()); err != nil {
		txn.Rollback()
		if key, _, ok := storage.CheckErrorDuplicate(err); ok {
			switch key {
			case "uid":
				return domain.ErrDuplicateImageID
			}
		}
		return err
	}

	if s.staticRepository == nil {
		return txn.Commit()
	}

	if err := s.staticRepository.Upload(storage.Image, image.ID(), image.Data()); err != nil {
		txn.Rollback()
		return domain.ErrFailedToUpload
	}

	return txn.Commit()
}

func (s *ImageStorage) Remove(image *model.Image) error {
	stmt := s.db.MustPrepare(`DELETE FROM image WHERE uid = ?`)
	defer stmt.Close()

	txn, err := s.db.Begin()
	if err != nil {
		return err
	}

	res, err := stmt.Exec(image.ID())
	if err != nil {
		txn.Rollback()
		return err
	}

	if rowsAffected, e := res.RowsAffected(); e != nil {
		txn.Rollback()
		return e
	} else if rowsAffected == 0 {
		txn.Rollback()
		return domain.ErrNoSuchImage
	}

	if s.staticRepository != nil {
		if err := s.staticRepository.Delete(storage.Image, image.ID()); err != nil {
			txn.Rollback()
			return err
		}
	}

	return txn.Commit()
}

func (s *ImageStorage) FindByID(id string) (*model.Image, error) {
	stmt := s.db.MustPrepare(`SELECT uid, user, created_at FROM image WHERE uid = ?`)
	defer stmt.Close()

	var (
		imageID         string
		imageUploaderID uint64
		imageUploadedAt time.Time
	)
	if err := stmt.QueryRow(id).Scan(&imageID, &imageUploaderID, &imageUploadedAt); err != nil {
		switch err {
		case storage.ErrNoRows:
			return nil, domain.ErrNoSuchImage
		}
	}

	return model.NewImage(imageID, imageUploaderID, imageUploadedAt), nil
}
