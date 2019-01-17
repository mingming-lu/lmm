package application

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lmm/api/service/asset/application/command"
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
	imageService    service.ImageService
	imageEncoder    service.ImageEncoder
	assetFinder     service.AssetFinder
}

// NewService creates a new image application service
func NewService(
	assetFinder service.AssetFinder,
	assetRepository repository.AssetRepository,
	imageService service.ImageService,
	imageEncoder service.ImageEncoder,
	uploaderService service.UploaderService,
) *Service {
	return &Service{
		assetFinder:     assetFinder,
		assetRepository: assetRepository,
		imageService:    imageService,
		imageEncoder:    imageEncoder,
		uploaderService: uploaderService,
	}
}

// UploadAsset handles upload asset command
func (app *Service) UploadAsset(c context.Context, cmd *command.UploadAsset) error {
	uploader, err := app.uploaderService.FromUserID(c, cmd.UserID())
	if err != nil {
		return errors.Wrap(err, cmd.UserID())
	}

	t := cmd.Type()
	switch t {
	case model.Image, model.Photo:
		return app.uploadImage(c, uploader, t, cmd.Data())
	default:
		return errors.Wrap(domain.ErrUnsupportedAssetType, t.String())
	}
}

func (app *Service) uploadImage(c context.Context, uploader *model.Uploader, assetType model.AssetType, data []byte) error {
	dst, ext, err := app.imageEncoder.Encode(c, data)
	if err != nil {
		if err == domain.ErrUnsupportedImageFormat {
			return errors.Wrap(err, ext)
		}
		return err
	}

	name := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.New(), dst).String()))
	asset := model.NewAsset(assetType, name+"."+ext, uploader, dst)

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

func (app *Service) SetPhotoAlternateTexts(c context.Context, cmd *command.SetImageAlternateTexts) error {
	asset, err := app.assetRepository.FindByName(c, cmd.ImageName())
	if err != nil {
		return errors.Wrap(domain.ErrNoSuchAsset, err.Error())
	}

	if asset.Type() != model.Photo {
		return domain.ErrInvalidTypeNotAPhoto
	}

	alts := make([]*model.Alt, len(cmd.AltNames()))
	for i, name := range cmd.AltNames() {
		alts[i] = model.NewAlt(asset.Name(), name)
	}

	return app.imageService.SetAlt(c, asset, alts)
}
