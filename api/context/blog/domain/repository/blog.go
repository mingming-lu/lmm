package repository

import "lmm/api/context/blog/domain/model"

type BlogRepository interface {
	Add(blog *model.Blog) error
	Update(blog *model.Blog) error
	FindAll(count, page int) ([]*model.Blog, int, error)
	FindAllByCategory(category *model.Category, count, page int) ([]*model.Blog, int, error)
	FindByID(id uint64) (*model.Blog, error)
	SetBlogCategory(blog *model.Blog, category *model.Category) error
	RemoveBlogCategory(blog *model.Blog) error
}
