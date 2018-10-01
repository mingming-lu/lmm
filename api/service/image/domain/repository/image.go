package repository

import (
	"context"

	"lmm/api/service/image/domain/model"
)

// ImageRepository provides a interface to deal with persistence of image
type ImageRepository interface {
	Save(c context.Context, image *model.Image) error
	Remove(c context.Context, image *model.Image) error
}
