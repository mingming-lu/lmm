package application

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"

	"lmm/api/service/asset/domain/model"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
)

// Service struct
type Service struct {
	uploaderService service.UploaderService
	assetRepository repository.AssetRepository
}

// NewService creates a new image application service
func NewService(
	assetRepository repository.AssetRepository,
	uploaderService service.UploaderService,
) *Service {
	return &Service{
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
