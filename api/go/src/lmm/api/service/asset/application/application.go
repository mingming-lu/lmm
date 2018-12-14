package application

import (
	"context"
	"encoding/base64"
	"image"
	"mime/multipart"

	// image encoder
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lmm/api/service/asset/domain"
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
func (app *Service) UploadImage(c context.Context, username, extention string, file multipart.File) error {
	return app.uploadAsset(c, model.Image, username, extention, file)
}

// UploadPhoto uploads photo
func (app *Service) UploadPhoto(c context.Context, username, extention string, file multipart.File) error {
	return app.uploadAsset(c, model.Photo, username, extention, file)
}

func (app *Service) uploadAsset(c context.Context, assetType model.AssetType,
	username string,
	extention string,
	file multipart.File,
) error {
	uploader, err := app.uploaderService.FromUserName(c, username)
	if err != nil {
		return err
	}

	switch assetType {
	case model.Image, model.Photo:
		return app.uploadImage(c, uploader, assetType, extention, file)
	default:
		return domain.ErrUnsupportedAssetType
	}
}

func (app *Service) uploadImage(c context.Context, uploader *model.Uploader, assetType model.AssetType, extention string, file multipart.File) error {
	src, _, err := image.Decode(file)
	if err != nil {
		return errors.Wrap(err, "invalid multipart.File, not an image")
	}

	dst, err := service.DefaultImageEncoder.Encode(c, src)
	name := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.New(), dst).String()))
	asset := model.NewAsset(assetType, name+".jpeg", uploader, dst)

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
