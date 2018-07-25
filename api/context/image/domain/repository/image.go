package repository

import "lmm/api/context/image/domain/model"

type ImageRepository interface {
	Add(*model.Image) error
	Remove(*model.Image) error
	FindByID(id string) (*model.Image, error)
}
