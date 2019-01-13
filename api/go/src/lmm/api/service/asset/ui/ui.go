package ui

import (
	"fmt"
	"mime/multipart"

	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/asset/application"
	"lmm/api/service/asset/application/command"
	"lmm/api/service/asset/domain"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
)

const (
	maxFormDataSize = 10 << 20 // 10MB
	maxImageSize    = 2 << 20  // 2MB
)

var (
	errRequestBodyTooLarge  = errors.New("request body too large")
	errImageMaxSizeExceeded = errors.New("the max size of image to upload cannot be larger than 2MB")
	errMaxUploadExcceed     = errors.New("only one image can be uploaded every time")
	errNoImageToUpload      = errors.New("please upload an image")
	errNotAllowedImageType  = errors.New("only gif, jpeg and png are allowed to upload")
)

// UI is a image UI
type UI struct {
	appService *application.Service
}

// New creates a new UI pointer
func New(
	assetFinder service.AssetFinder,
	assetRepository repository.AssetRepository,
	imageService service.ImageService,
	userAdapter service.UploaderService,
) *UI {
	appService := application.NewService(
		assetFinder,
		assetRepository,
		imageService,
		userAdapter,
	)
	return &UI{appService: appService}
}

// UploadImage handles POST /v1/assets/images
func (ui *UI) UploadImage(c http.Context) {
	ui.uploadImage(c, "image")
}

// UploadPhoto handles POST /v1/assets/photos
func (ui *UI) UploadPhoto(c http.Context) {
	ui.uploadImage(c, "photo")
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

func (ui *UI) uploadImage(c http.Context, keyName string) {
	userID := c.Request().Header.Get("X-LMM-ID")
	if userID == "" {
		http.Unauthorized(c)
		return
	}

	file, err := formImageData(c, keyName)
	if err != nil {
		http.Log().Warn(c, err.Error())
		c.String(http.StatusBadRequest, errors.Cause(err).Error())
		return
	}
	defer file.Close()

	err = ui.appService.UploadAsset(c, command.NewUploadAsset(userID, keyName, file))
	switch errors.Cause(err) {
	case nil:
		c.String(http.StatusCreated, "uploaded")
	case domain.ErrUnsupportedImageFormat:
		http.Log().Warn(c, err.Error())
		c.String(http.StatusBadRequest, domain.ErrUnsupportedImageFormat.Error())
	case domain.ErrNoSuchUser:
		http.Log().Warn(c, err.Error())
		http.Unauthorized(c)
	case domain.ErrUnsupportedAssetType:
		// unexpected (if no bug)
		http.Log().Panic(c, err.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}

func formImageData(c http.Context, imageKey string) (multipart.File, error) {
	if err := c.Request().ParseMultipartForm(maxFormDataSize); err != nil {
		return nil, errors.Wrap(errRequestBodyTooLarge, err.Error())
	}
	assets := c.Request().MultipartForm.File[imageKey]
	switch len(assets) {
	case 0:
		return nil, errNoImageToUpload
	case 1:
	default:
		return nil, errors.Wrap(errMaxUploadExcceed, fmt.Sprintf("attempt to upload %d assets", len(assets)))
	}

	asset := assets[0]
	if asset.Size > maxImageSize {
		return nil, errImageMaxSizeExceeded
	}

	// check type
	contentType := asset.Header.Get("Content-Type")
	switch contentType {
	case "image/bmp":
	case "image/gif":
	case "image/jpeg":
	case "image/png":
	case "image/webp":
	default:
		return nil, errors.Wrap(errNotAllowedImageType, contentType)
	}

	f, err := asset.Open() // must open
	if err != nil {
		http.Log().Panic(c, err.Error())
	}

	return f, nil
}
