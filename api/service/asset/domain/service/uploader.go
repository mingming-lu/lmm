package service

import (
	"context"

	"lmm/api/service/image/domain/model"
)

// UploaderService provides interface to convert user to uploader
type UploaderService interface {
	FromUserName(c context.Context, name string) (*model.Uploader, error)
}
