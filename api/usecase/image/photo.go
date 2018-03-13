package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
)

func AllPhotos(userID int64) ([]model.Minimal, error) {
	return repo.FetchAllPhotos(userID)
}

func ShowOnPhotos(userID, imageID int64) error {
	// TODO check image id
	return repo.MarkAsPhoto(userID, imageID)
}
