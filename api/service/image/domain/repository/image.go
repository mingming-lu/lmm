package repository

import (
	"lmm/api/service/image/domain"
	"lmm/api/service/image/domain/model"
)

type ImageRepository interface {
	Add(*model.ImageWithData) error
	Remove(*model.Image) error
	FindByID(id string) (*model.Image, error)
	FindByType(domain.ImageType, int, int) ([]*model.Image, bool, error)
	Find(count, page int) ([]*model.Image, bool, error)
	MarkAs(*model.Image, domain.ImageType) error
}
