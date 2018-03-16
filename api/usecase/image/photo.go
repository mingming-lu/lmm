package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
)

func AllPhotos(userID int64) ([]model.Minimal, error) {
	return repo.FetchAllPhotos(userID)
}

func TurnOnPhotoSwitch(userID int64, imageName string) error {
	return togglePhotoSwitch(userID, imageName, true)
}

func TurnOffPhotoSwitch(userID int64, imageName string) error {
	return togglePhotoSwitch(userID, imageName, false)
}

func togglePhotoSwitch(userID int64, imageName string, shown bool) error {
	image, err := repo.ByName(userID, imageName)
	if err != nil {
		return err
	}
	return repo.SavePhoto(userID, image.ID, shown)
}
