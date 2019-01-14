package service

import (
	"context"
	"lmm/api/service/asset/domain/model"
)

type ImageService interface {
	SetAlt(context.Context, *model.AssetDescriptor, []*model.Alt) error
}
