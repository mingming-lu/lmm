package persistence

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"lmm/api/clock"
	"lmm/api/service/asset/domain/model"
	"lmm/api/storage/db"
	"lmm/api/storage/uploader"
)

// AssetStorage s a AssetRepository implementation
type AssetStorage struct {
	db       db.DB
	uploader uploader.Uploader
}

// NewAssetStorage creates new AssetStorage
func NewAssetStorage(db db.DB, uploader uploader.Uploader) *AssetStorage {
	return &AssetStorage{db: db, uploader: uploader}
}

type asset struct {
	id       uint
	typeCode uint8
	name     string
}

// FindAssetByName implementation
func (s *AssetStorage) FindAssetByName(c context.Context, name string) (*model.AssetDescriptor, error) {
	tx, err := s.db.Begin(c, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		panic(err)
	}

	a, err := s.findAssetByName(c, tx, name)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return nil, errors.Wrap(err, e.Error())
		}
		return nil, err
	}

	return model.NewAssetDescriptor(a.name, s.decodeAssetType(a.typeCode)), tx.Commit()
}

func (s *AssetStorage) findAssetByName(c context.Context, tx db.Tx, name string) (*asset, error) {
	stmt, err := tx.Prepare(c, `select id, name, type from asset where name = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	a := &asset{}
	if err := stmt.QueryRow(c, name).Scan(&a.id, &a.name, &a.typeCode); err != nil {
		return nil, err
	}

	return a, nil
}

// FindPhotoByName implementation
func (s *AssetStorage) FindPhotoByName(c context.Context, name string) (*model.PhotoDescriptor, error) {
	tx, err := s.db.Begin(c, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare(c, `select name from image_alt where asset = ?`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	a, err := s.findAssetByName(c, tx, name)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return nil, errors.Wrap(err, e.Error())
		}
		return nil, err
	}
	photo := model.NewPhotoDescriptor(a.id, a.name)

	rows, err := stmt.Query(c, a.id)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return nil, errors.Wrap(err, e.Error())
		}
		return nil, err
	}

	var altName string
	for rows.Next() {
		if err := rows.Scan(&altName); err != nil {
			panic(err)
		}
		if err := photo.AddAlternateText(altName); err != nil {
			if e := tx.Rollback(); e != nil {
				return nil, errors.Wrap(err, e.Error())
			}
			return nil, err
		}
	}
	rows.Close()

	return photo, tx.Commit()
}

// Save implementation
func (s *AssetStorage) Save(c context.Context, asset *model.Asset) error {
	tx, err := s.db.Begin(c, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(c, `insert into asset (name, type, user, created_at) values (?, ?, ?, ?)`)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	if _, err := stmt.Exec(c,
		asset.Name(),
		s.encodeAssetType(asset.Type()),
		asset.Uploader().ID(),
		clock.Now(),
	); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	if err := s.uploader.Upload(c,
		asset.Name(),
		asset.Data(),
		uploader.ImageUploaderConfig{
			Type: asset.Type().String(),
		},
	); err != nil {
		if e := tx.Rollback(); e != nil {
			return errors.Wrap(err, e.Error())
		}
		return err
	}

	return tx.Commit()
}

// Remove implementation
func (s *AssetStorage) Remove(c context.Context, asset *model.Asset) error {
	panic("not implemented")
}

func (s *AssetStorage) encodeAssetType(assetType model.AssetType) uint8 {
	switch assetType.String() {
	case "image":
		return 0
	case "photo":
		return 1
	default:
		panic("invalid asset type: '" + assetType.String() + "'")
	}
}

func (s *AssetStorage) decodeAssetType(code uint8) model.AssetType {
	switch code {
	case 0:
		return model.Image
	case 1:
		return model.Photo
	default:
		panic("invalid asset code: '" + string(code) + "'")
	}
}
