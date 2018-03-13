package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
)

func AllPhotos(userID int64) ([]model.Minimal, error) {
	return repo.FetchAllPhotos(userID)
}

func ShowOnPhotos(userID int64, imageName string) error {
	image, err := repo.ByName(userID, imageName)
	if err != nil {
		return err
	}
	return repo.MarkAsPhoto(userID, image.ID)
}
