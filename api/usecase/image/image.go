package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
	"lmm/api/domain/service/base64"
	"lmm/api/domain/service/uuid"
)

func Upload(userID uint64, bulkData [][]byte) error {
	imageData := make([]model.ImageData, 0)
	for _, image := range bulkData {
		imgData := model.ImageData{
			Name: base64.Encode([]byte(uuid.New())),
			Data: image,
		}
		imageData = append(imageData, imgData)
	}
	return repo.Add(userID, imageData)
}

func AllImages(userID uint64) ([]model.Minimal, error) {
	return repo.FetchAllImage(userID)
}
