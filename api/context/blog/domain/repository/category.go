package repository

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/domain/repository"
)

type CategoryRepository interface {
	repository.Repository
	Add(category *model.Category) error
	Update(categoryRepo *model.Category) error
	FindAll(count, page int) ([]*model.Category, error)
}

type categoryRepo struct {
	repository.Default
}

func NewCategoryRepository() CategoryRepository {
	return new(categoryRepo)
}

func (c *categoryRepo) Add(category *model.Category) error {
	return nil
}

func (c *categoryRepo) Update(category *model.Category) error {
	return nil
}

func (c *categoryRepo) FindAll(count, page int) ([]*model.Category, error) {
	return nil, nil
}
