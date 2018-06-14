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
	db := c.DB()
	defer db.Close()

	stmt := db.MustPrepare(`INSERT INTO category (id, name) VALUES (?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(category.ID(), category.Name())
	return err
}

func (c *categoryRepo) Update(category *model.Category) error {
	panic("not implemented")
	return nil
}

func (c *categoryRepo) FindAll(count, page int) ([]*model.Category, error) {
	panic("not implemented")
	return nil, nil
}
