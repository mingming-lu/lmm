package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
	"lmm/api/domain/service/uuid"
)

func Upload(userID int64, t model.ImageType) error {
	name := uuid.New()
	return repo.Add(userID, t, model.BaseURL+name)
}

func Find(userID int64, t model.ImageType) ([]model.Image, error) {
	return repo.Fetch(userID, t)
}
