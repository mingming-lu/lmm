package ui

import (
	"io/ioutil"
	"mime/multipart"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"lmm/api/http"
	userModel "lmm/api/service/auth/domain/model"
	"lmm/api/service/image/application"
	"lmm/api/service/image/domain/repository"
	"lmm/api/service/image/domain/service"
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
func New(imageRepository repository.ImageRepository, userAdapter service.UploaderService) *UI {
	appService := application.NewService(imageRepository, userAdapter)
	return &UI{appService: appService}
}

// Upload handles POST /v1/images
func (ui *UI) Upload(c http.Context) {
	user, ok := c.Value(http.StrCtxKey("user")).(*userModel.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	data, ext, err := formImageData(c, "image")
	if err != nil {
		zap.L().Warn(err.Error(), zap.String("request_id", c.Request().Header.Get("X-Request-ID")))
		c.String(http.StatusBadRequest, errors.Cause(err).Error())
		return
	}

	// upload
	if err := ui.appService.UploadImage(c, user.Name(), data, ext); err != nil {
		panic(err)
	}

	c.String(http.StatusCreated, "upload")
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

// ListImages handles GET /v1/assets/images
func (ui *UI) ListImages() {
}

// UploadPhoto handles POST /v1/assets/photos
func (ui *UI) UploadPhoto(c http.Context) {
}

// ListPhotos handles GET /v1/assets/photos
func (ui *UI) ListPhotos(c http.Context) {
}
