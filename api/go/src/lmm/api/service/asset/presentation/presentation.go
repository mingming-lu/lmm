package presentation

import (
	"lmm/api/http"
	"lmm/api/pkg/auth"
	"lmm/api/service/asset/usecase"

	"github.com/pkg/errors"
)

var (
	ErrUnsupportType = errors.New("unsupport type")
)

type Presentation struct {
	usecase *usecase.Usecase
}

func New(app *usecase.Usecase) *Presentation {
	return &Presentation{usecase: app}
}

// PostV1Photos handles POST /v1/photos
func (p *Presentation) PostV1Photos(c http.Context) {
	user, ok := auth.FromContext(c)
	if !ok {
		http.Unauthorized(c)
		return
	}

	f, fh, err := c.Request().FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			c.String(http.StatusBadRequest, "photo required")
			return
		}
		http.Log().Warn(c, err.Error())
		http.BadRequest(c)
		return
	}

	contentType := fh.Header.Get("Content-Type")
	switch contentType {
	case "image/gif", "image/jpeg", "image/jpg", "image/png":
	default:
		c.String(http.StatusBadRequest, ErrUnsupportType.Error())
	}

	url, err := p.usecase.UploadPhoto(c, &usecase.AssetToUpload{
		ContentType: contentType,
		DataSource:  f,
		Filename:    fh.Filename,
		UserID:      user.ID,
	})

	switch err {
	case nil:
		c.Header("Location", url)
		http.Created(c)
	default:
		http.Log().Panic(c, err.Error())
	}
}

type photoList struct {
	Items      []*usecase.Photo `json:"items"`
	NextCursor string           `json:"next_cursor"`
}

// GetV1Photos handles GET /v1/photos
func (p *Presentation) GetV1Photos(c http.Context) {
	photos, cursor, err := p.usecase.ListPhotos(c,
		c.Request().QueryParamOrDefault("count", "10"),
		c.Request().QueryParamOrDefault("cursor", ""),
	)

	if err != nil {
		http.Log().Warn(c, err.Error())
		http.NotFound(c)
		return
	}

	c.JSON(http.StatusOK, &photoList{
		Items:      photos,
		NextCursor: cursor,
	})
}

// PatchV1Photos handles PATCH /v1/photos/:photo
func (p *Presentation) PatchV1Photos(c http.Context) {
}

// PostV1Assets handles POST /v1/assets
// This endpoint is to upload common assets
func (p *Presentation) PostV1Assets(c http.Context) {
}
