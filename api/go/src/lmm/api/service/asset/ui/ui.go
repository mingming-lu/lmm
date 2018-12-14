package ui

import (
	"context"
	"lmm/api/service/asset/domain"
	"mime/multipart"

	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/asset/application"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
)

const (
	maxFormDataSize = 32 << 20 // 32MB
	maxImageSize    = 2 << 20  // 2MB
)

var (
	errImageMaxSizeExceeded = errors.New("the max size of image to upload cannot be larger than 2MB")
	errNotAllowedImageType  = errors.New("only gif, jpeg and png are allowed to upload")
	errMaxUploadExcceed     = errors.New("only one image can be uploaded every time")
	errNoImageToUpload      = errors.New("please upload an image")
)

// UI is a image UI
type UI struct {
	appService *application.Service
}

// New creates a new UI pointer
func New(
	assetFinder service.AssetFinder,
	assetRepository repository.AssetRepository,
	userAdapter service.UploaderService,
) *UI {
	appService := application.NewService(assetFinder, assetRepository, userAdapter)
	return &UI{appService: appService}
}

// UploadImage handles POST /v1/assets/images
func (ui *UI) UploadImage(c http.Context) {
	ui.upload(c, "image")
}

// UploadPhoto handles POST /v1/assets/photos
func (ui *UI) UploadPhoto(c http.Context) {
	ui.upload(c, "photo")
}

// ListImages handles GET /v1/assets/images
func (ui *UI) ListImages(c http.Context) {
	collection, err := ui.appService.ListImages(c,
		c.Request().QueryParamOrDefault("page", "1"),
		c.Request().QueryParamOrDefault("perPage", "100"),
	)
	if err == nil {
		c.JSON(http.StatusOK, imageCollectionToJSON(collection))
		return
	}

	switch errors.Cause(err) {
	case application.ErrInvalidPage, application.ErrInvalidPerPage:
		http.NotFound(c)
	default:
		http.Log().Panic(c, err.Error())
	}
}

// ListPhotos handles GET /v1/assets/photos
func (ui *UI) ListPhotos(c http.Context) {
	collection, err := ui.appService.ListPhotos(c,
		c.Request().QueryParamOrDefault("page", "1"),
		c.Request().QueryParamOrDefault("perPage", "10"),
	)
	if err == nil {
		c.JSON(http.StatusOK, photoCollectionToJSON(collection))
		return
	}

	switch errors.Cause(err) {
	case application.ErrInvalidPage, application.ErrInvalidPerPage:
		http.NotFound(c)
	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) upload(c http.Context, keyName string) {
	userName := c.Request().Header.Get("X-LMM-ID")
	if userName == "" {
		http.Unauthorized(c)
		return
	}

	file, ext, err := formImageData(c, keyName)
	if err != nil {
		http.Log().Error(c, err.Error())
		c.String(http.StatusBadRequest, errors.Cause(err).Error())
		return
	}

	// upload
	switch keyName {
	case "image":
		if err := ui.appService.UploadImage(c, userName, ext, file); err != nil {
			panic(err)
		}
	case "photo":
		if err := ui.appService.UploadPhoto(c, userName, ext, file); err != nil {
			panic(err)
		}
	default:
		panic("unknown key name: '" + keyName + "'")
	}

	c.String(http.StatusCreated, "uploaded")
}

func formImageData(c http.Context, imageKey string) (multipart.File, string, error) {
	if err := c.Request().ParseMultipartForm(maxFormDataSize); err != nil {
		http.Log().Error(c, err.Error())
		return nil, "", errors.Wrap(errImageMaxSizeExceeded, err.Error())
	}
	assets := c.Request().MultipartForm.File[imageKey]

	if len(assets) == 0 {
		return nil, "", errNoImageToUpload
	}

	if len(assets) > 1 {
		http.Log().Error(c, "attend to upload multiple images")
		return nil, "", errMaxUploadExcceed
	}

	return openImage(c, assets[0])
}

func openImage(c context.Context, fh *multipart.FileHeader) (multipart.File, string, error) {
	// check type
	ext := ""
	contentType := fh.Header.Get("Content-Type")
	switch contentType {
	case "image/gif":
		ext = domain.ImageExtensionGIF
	case "image/jpeg":
		ext = domain.ImageExtensionJPEG
	case "image/png":
		ext = domain.ImageExtensionPNG
	default:
		return nil, "", errNotAllowedImageType
	}

	f, err := fh.Open() // must open
	if err != nil {
		http.Log().Panic(c, err.Error())
	}
	defer f.Close()

	return f, ext, nil
}
