package factory

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/domain/factory"
)

func NewTag(blogID uint64, name string) (*model.Tag, error) {
	id, err := factory.Default().GenerateID()
	if err != nil {
		return nil, err
	}
	return model.NewTag(id, blogID, name)
}
