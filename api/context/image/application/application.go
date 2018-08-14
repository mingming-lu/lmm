package application

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/image/domain"
	"lmm/api/context/image/domain/factory"
	"lmm/api/context/image/domain/model"
	"lmm/api/context/image/domain/repository"
	"lmm/api/utils/strings"
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

func (app *AppService) FetchImages(countStr, pageStr string) ([]*model.Image, bool, error) {
	if countStr == "" {
		countStr = "100"
	}
	if pageStr == "" {
		pageStr = "1"
	}
	count, err := strings.StrToInt(countStr)
	if err != nil {
		return nil, false, domain.ErrInvalidCount
	}
	if count < 0 {
		return nil, false, domain.ErrInvalidCount
	}
	page, err := strings.StrToInt(pageStr)
	if err != nil {
		return nil, false, domain.ErrInvalidPage
	}
	if page < 1 {
		return nil, false, domain.ErrInvalidPage
	}
	return app.imageRepo.Find(count, page)
}
