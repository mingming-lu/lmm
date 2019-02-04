package service

import (
	"context"

	"lmm/api/service/asset/domain"
	"lmm/api/service/asset/domain/model"
	"lmm/api/storage/db"

	"github.com/pkg/errors"
)

type ImageService struct {
	db db.DB
}

func NewImageService(db db.DB) *ImageService {
	return &ImageService{db: db}
}

func (s *ImageService) SetAlt(c context.Context, asset *model.AssetDescriptor, alts []*model.Alt) error {
	tx, err := s.db.Begin(c, nil)
	if err != nil {
		return err
	}

	selectAsset := tx.Prepare(c, "select id from asset where name = ?")
	if err != nil {
		return db.RollbackWithError(tx, err)
	}
	defer selectAsset.Close()

	var assetID uint
	if err := selectAsset.QueryRow(c, asset.Name()).Scan(&assetID); err != nil {
		err = errors.Wrap(domain.ErrNoSuchAsset, err.Error())
		return db.RollbackWithError(tx, err)
	}

	deleteAltByAssetID := tx.Prepare(c, "delete from image_alt where asset = ?")
	if err != nil {
		return db.RollbackWithError(tx, err)
	}
	defer deleteAltByAssetID.Close()

	if _, err := deleteAltByAssetID.Exec(c, assetID); err != nil && err != db.ErrNoRows {
		return db.RollbackWithError(tx, err)
	}

	insertAlt := tx.Prepare(c, "insert into image_alt (asset, name) values(?, ?)")
	if err != nil {
		return db.RollbackWithError(tx, err)
	}
	defer insertAlt.Close()

	for _, alt := range alts {
		if _, err := insertAlt.Exec(c, assetID, alt.Name()); err != nil {
			return db.RollbackWithError(tx, err)
		}
	}

	return tx.Commit()
}
