package ui

import (
	"bytes"
	"fmt"
	"mime/multipart"
	net "net/http"

	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/asset/application"
	"lmm/api/service/asset/application/command"
	"lmm/api/service/asset/domain"
	"lmm/api/service/asset/domain/repository"
	"lmm/api/service/asset/domain/service"
)

var (
	maxFormDataSize int64 = 10 << 20 // 10MB
	maxImageSize    int64 = 2 << 20  // 2MB
)

var (
	errImageMaxSizeExceeded = errors.New("the max size of image to upload cannot be larger than 2MB")
	errMaxUploadExcceed     = errors.New("only one image can be uploaded every time")
	errNoImageToUpload      = errors.New("please upload an image")
	errNotAllowedImageType  = errors.New("unsupported image format")
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
	imageEncoder service.ImageEncoder,
	userAdapter service.UploaderService,
) *UI {
	appService := application.NewService(
		assetFinder,
		assetRepository,
		imageService,
		imageEncoder,
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

	buf, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		panic("expected a *bytes.Buffer")
	}
	defer bufPool.Put(buf)

	err := formImageData(c, keyName, buf)
	switch errors.Cause(err) {
	case nil:
	case multipart.ErrMessageTooLarge:
		c.String(http.StatusBadRequest, err.Error())
		return
	case http.ErrNotMultipart:
		http.BadRequest(c)
		return
	default:
		http.Log().Warn(c, err.Error())
		c.String(http.StatusBadRequest, errors.Cause(err).Error())
		return
	}

	err = ui.appService.UploadAsset(c, command.NewUploadAsset(userID, keyName, buf.Bytes()))
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

func formImageData(c http.Context, imageKey string, buf *bytes.Buffer) error {
	if err := c.Request().ParseMultipartForm(maxFormDataSize); err != nil {
		return err
	}
	assets := c.Request().MultipartForm.File[imageKey]
	switch len(assets) {
	case 0:
		return errNoImageToUpload
	case 1:
	default:
		return errors.Wrap(errMaxUploadExcceed, fmt.Sprintf("attempt to upload %d assets", len(assets)))
	}

	asset := assets[0]
	if asset.Size > maxImageSize {
		return errImageMaxSizeExceeded
	}

	f, err := asset.Open() // must open
	if err != nil {
		http.Log().Panic(c, err.Error())
	}
	defer f.Close()

	buf.Reset()
	if size, err := buf.ReadFrom(f); err != nil {
		panic(err)
	} else if size > maxImageSize {
		return errImageMaxSizeExceeded
	}

	contentType := net.DetectContentType(buf.Bytes())
	switch contentType {
	case "image/bmp":
	case "image/gif":
	case "image/jpeg":
	case "image/png":
	case "image/webp":
	default:
		return errors.Wrap(errNotAllowedImageType, contentType)
	}

	return nil
}

// PutPhotoAlternateTexts handles PUT /v1/assets/photos/:photo/alts
func (ui *UI) PutPhotoAlternateTexts(c http.Context) {
	userID := c.Request().Header.Get("X-LMM-ID")
	if userID == "" {
		http.Unauthorized(c)
		return
	}

	imageName := c.Request().PathParam("photo")

	altNames := putPhotoAltsRequestBody{}
	if err := c.Request().Bind(&altNames); err != nil {
		http.Log().Warn(c, err.Error())
		http.BadRequest(c)
		return
	}

	err := ui.appService.SetPhotoAlternateTexts(c,
		command.NewSetImageAlternateTexts(imageName, altNames.Names),
	)

	switch errors.Cause(err) {
	case nil:
		http.NoContent(c)
	case domain.ErrInvalidTypeNotAPhoto:
		c.String(http.StatusBadRequest, domain.ErrInvalidTypeNotAPhoto.Error())
	case domain.ErrNoSuchAsset:
		c.String(http.StatusNotFound, domain.ErrNoSuchAsset.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}

// GetPhotoDescription handles GET /v1/assets/photos/:photo
func (ui *UI) GetPhotoDescription(c http.Context) {
	userID := c.Request().Header.Get("X-LMM-ID")
	if userID == "" {
		http.Unauthorized(c)
		return
	}

	descriptor, err := ui.appService.GetPhotoDescription(c, c.Request().PathParam("photo"))
	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, &photoListItem{
			Name: descriptor.Name(),
			Alts: descriptor.AlternateTexts(),
		})
	case domain.ErrNoSuchPhoto:
		http.Log().Warn(c, err.Error())
		c.String(http.StatusNotFound, domain.ErrNoSuchPhoto.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}
