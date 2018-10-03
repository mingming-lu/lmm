package ui

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"lmm/api/http"
	"lmm/api/service/asset/application"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
	userModel "lmm/api/service/auth/domain/model"
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
		c.Request().QueryParam("limit"),
		c.Request().QueryParam("nextCursor"),
	)
	if err == nil {
		c.JSON(http.StatusOK, imageCollectionToJSON(collection))
		return
	}

	zap.L().Warn(err.Error(), zap.String("request_id", c.Request().Header.Get("X-Request-ID")))

	err = errors.Cause(err)
	switch err {
	case application.ErrInvalidCursor, application.ErrInvalidLimit:
		c.String(http.StatusBadRequest, err.Error())
	default:
		panic(err)
	}
}

// ListPhotos handles GET /v1/assets/photos
func (ui *UI) ListPhotos(c http.Context) {
}

func (ui *UI) upload(c http.Context, keyName string) {
	user, ok := c.Value(http.StrCtxKey("user")).(*userModel.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	data, ext, err := formImageData(c, keyName)
	if err != nil {
		zap.L().Warn(err.Error(), zap.String("request_id", c.Request().Header.Get("X-Request-ID")))
		c.String(http.StatusBadRequest, errors.Cause(err).Error())
		return
	}

	// upload
	switch keyName {
	case "image":
		if err := ui.appService.UploadImage(c, user.Name(), data, ext); err != nil {
			panic(err)
		}
	case "photo":
		if err := ui.appService.UploadPhoto(c, user.Name(), data, ext); err != nil {
			panic(err)
		}
	default:
		panic("unknown key name: '" + keyName + "'")
	}

	c.String(http.StatusCreated, "uploaded")
}

func formImageData(c http.Context, imageKey string) ([]byte, string, error) {
	if err := c.Request().ParseMultipartForm(maxFormDataSize); err != nil {
		zap.L().Warn(err.Error(), zap.String("request_id", c.Request().Header.Get("X-Request-ID")))
		return nil, "", errors.Wrap(errImageMaxSizeExceeded, err.Error())
	}
	assets := c.Request().MultipartForm.File[imageKey]

	if len(assets) == 0 {
		return nil, "", errNoImageToUpload
	}

	if len(assets) > 1 {
		zap.L().Warn("attend to upload multiple images", zap.String("request_id", c.Request().Header.Get("X-Request-ID")))
		return nil, "", errMaxUploadExcceed
	}

	return openImage(assets[0])
}

func openImage(fh *multipart.FileHeader) ([]byte, string, error) {
	// check type
	ext := ""
	contentType := fh.Header.Get("Content-Type")
	switch contentType {
	case "image/gif":
		ext = "gif"
	case "image/jpeg":
		ext = "jpeg"
	case "image/png":
		ext = "png"
	default:
		return nil, "", errNotAllowedImageType
	}

	f, err := fh.Open()
	if err != nil {
		// must open file otherwise go die
		panic(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		// if it can't read all, just go die
		panic(err)
	}
	return data, ext, nil
}
