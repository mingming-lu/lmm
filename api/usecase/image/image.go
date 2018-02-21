package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
	"lmm/api/domain/service/base64"
	"lmm/api/domain/service/uuid"
)

func Upload(userID int64, imageTypeStr string, data []byte) error {
	imageType := toImageType(imageTypeStr)
	name := base64.Encode([]byte(uuid.New()))
	return repo.Add(userID, imageType, name, data)
}

func Find(userID int64, t model.ImageType) ([]model.Image, error) {
	return repo.Fetch(userID, t)
}

func toImageType(imageType string) model.ImageType {
	switch imageType {
	case "avatar":
		return model.TypeAvatar
	case "blog":
		return model.TypeBlog
	case "photo":
		return model.TypePhoto
	default:
		panic("Unexpected image type: " + imageType)
	}
}
