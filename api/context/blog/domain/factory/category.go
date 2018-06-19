package factory

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/domain/factory"
)

func NewCategory(name string) (*model.Category, error) {
	id, err := factory.Default().GenerateID()
	if err != nil {
		return nil, err
	}
	return model.NewCategory(id, name)
}
