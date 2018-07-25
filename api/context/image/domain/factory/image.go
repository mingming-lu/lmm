package factory

import (
	"lmm/api/context/image/domain/model"
	"lmm/api/utils/uuid"
	"time"
)

func NewImage(userID uint64) *model.Image {
	id := uuid.New()
	createdAt := time.Now()

	return model.NewImage(id, userID, createdAt)
}
