package factory

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"lmm/api/service/image/domain/model"
)

func NewImage(userID uint64) *model.Image {
	id := uuid.New().String()
	hashedID := base64.RawURLEncoding.EncodeToString([]byte(id))
	createdAt := time.Now()

	return model.NewImage(hashedID, userID, createdAt)
}
