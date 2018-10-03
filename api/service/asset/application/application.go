package application

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lmm/api/service/asset/domain/model"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
	"lmm/api/util/stringutil"
)

var (
	ErrInvalidLimit  = errors.New("invalid limit")
	ErrInvalidCursor = errors.New("invalid cursor")
)

// Service struct
type Service struct {
	uploaderService service.UploaderService
	assetRepository repository.AssetRepository
	assetFinder     service.AssetFinder
}

// NewService creates a new image application service
func NewService(
	assetFinder service.AssetFinder,
	assetRepository repository.AssetRepository,
	uploaderService service.UploaderService,
) *Service {
	return &Service{
		assetFinder:     assetFinder,
		assetRepository: assetRepository,
		uploaderService: uploaderService,
	}
}

// UploadImage uploads image
func (app *Service) UploadImage(c context.Context, username string, data []byte, extention string) error {
	return app.uploadAsset(c, username, data, extention, model.Image)
}

// UploadPhoto uploads photo
func (app *Service) UploadPhoto(c context.Context, username string, data []byte, extention string) error {
	return app.uploadAsset(c, username, data, extention, model.Photo)
}

func (app *Service) uploadAsset(c context.Context,
	username string,
	data []byte,
	extention string,
	assetType model.AssetType,
) error {
	uploader, err := app.uploaderService.FromUserName(c, username)
	if err != nil {
		return err
	}

	name := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.New(), data).String()))
	asset := model.NewAsset(assetType, name+"."+extention, uploader, model.Data(data))

	return app.assetRepository.Save(c, asset)
}

// ListImages lists images by given limit and cursor
func (app *Service) ListImages(c context.Context, limitStr, nextCursorStr string) (*model.ImageCollection, error) {
	if limitStr == "" {
		limitStr = "30"
	}

	if nextCursorStr == "" {
		nextCursorStr = "0"
	}

	limit, err := stringutil.ParseUint(limitStr)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidLimit, err.Error())
	}

	nextCursor, err := stringutil.ParseUint(nextCursorStr)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidCursor, err.Error())
	}

	return app.assetFinder.FindAllImages(c, limit, nextCursor)
}

// ListPhotos lists images by given limit and cursor
func (app *Service) ListPhotos() (*model.PhotoCollection, error) {
	panic("not implemented")
}
