package model

import "lmm/api/domain/model"

type assetDescriptor struct {
	model.ValueObject
	name string
}

func (d *assetDescriptor) Name() string {
	return d.name
}
