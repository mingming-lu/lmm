package repository

import "lmm/api/context/blog/domain/model"

type CategoryRepository interface {
	Add(category *model.Category) error
	Update(categoryRepo *model.Category) error
	FindAll() ([]*model.Category, error)
	FindByID(id uint64) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	FindByBlog(blog *model.Blog) (*model.Category, error)
	Remove(category *model.Category) error
}
