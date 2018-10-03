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
	ErrInvalidPage    = errors.New("invalid page")
	ErrInvalidPerPage = errors.New("invalid perPage")
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

// ListImages lists images by given page and perPage
func (app *Service) ListImages(c context.Context, pageStr, perPageStr string) (*model.ImageCollection, error) {
	page, perPage, err := app.parseLimitAndCursorOrDefault(pageStr, perPageStr)
	if err != nil {
		return nil, err
	}

	return app.assetFinder.FindAllImages(c, page, perPage)
}

// ListPhotos lists images by given page and perPage
func (app *Service) ListPhotos(c context.Context, pageStr, perPageStr string) (*model.PhotoCollection, error) {
	page, perPage, err := app.parseLimitAndCursorOrDefault(pageStr, perPageStr)
	if err != nil {
		return nil, err
	}

	return app.assetFinder.FindAllPhotos(c, page, perPage)
}

func (app *Service) parseLimitAndCursorOrDefault(pageStr, perPageStr string) (uint, uint, error) {
	if pageStr == "" {
		pageStr = "1"
	}

	if perPageStr == "" {
		perPageStr = "30"
	}

	page, err := stringutil.ParseUint(pageStr)
	if err != nil {
		return 0, 0, errors.Wrap(ErrInvalidPage, err.Error())
	}
	if page < 1 {
		return 0, 0, errors.Wrap(ErrInvalidPage, "page can not be less than 1")
	}

	perPage, err := stringutil.ParseUint(perPageStr)
	if err != nil {
		return 0, 0, errors.Wrap(ErrInvalidPerPage, err.Error())
	}

	return page, perPage, nil
}
