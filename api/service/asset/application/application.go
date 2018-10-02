package application

import (
	"context"
	"encoding/base64"

	"lmm/api/service/image/domain/model"
	"lmm/api/service/image/domain/repository"
	"lmm/api/service/image/domain/service"

	"github.com/google/uuid"
)

// Service struct
type Service struct {
	uploaderService service.UploaderService
	imageRepository repository.ImageRepository
}

// NewService creates a new image application service
func NewService(
	imageRepository repository.ImageRepository,
	uploaderService service.UploaderService,
) *Service {
	return &Service{
		imageRepository: imageRepository,
		uploaderService: uploaderService,
	}
}

// UploadImage uploads image
func (app *Service) UploadImage(c context.Context, username string, data []byte, extention string) error {
	uploader, err := app.uploaderService.FromUserName(c, username)
	if err != nil {
		return err
	}

	name := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.New(), data).String()))

	image := model.NewImage(name+"."+extention, uploader, model.Data(data))

	return app.imageRepository.Save(c, image)
}
