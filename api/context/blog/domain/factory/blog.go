package factory

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/domain/factory"
	"time"
)

func NewBlog(userID uint64, title, text string) (*model.Blog, error) {
	id, err := factory.Default().GenerateID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return model.NewBlog(id, userID, title, text, now, now), nil
}
