package application

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/image/domain/factory"
	"lmm/api/context/image/domain/repository"
)

type AppService struct {
	imageRepo repository.ImageRepository
}

func NewAppService(imageRepo repository.ImageRepository) *AppService {
	return &AppService{
		imageRepo: imageRepo,
	}
}

func (app *AppService) UploadImage(user *account.User, data []byte) error {
	image := factory.NewImage(user.ID())
	return app.imageRepo.Add(image.WrapData(data))
}
