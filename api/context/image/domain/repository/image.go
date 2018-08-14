package repository

import (
	"lmm/api/context/image/domain"
	"lmm/api/context/image/domain/model"
)

type ImageRepository interface {
	Add(*model.ImageWithData) error
	Remove(*model.Image) error
	FindByID(id string) (*model.Image, error)
	Find(count, page int) ([]*model.Image, bool, error)
	MarkAs(*model.Image, domain.ImageType) error
}
