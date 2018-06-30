package repository

import "lmm/api/context/blog/domain/model"

type TagRepository interface {
	Add(tag *model.Tag) error
	FindByID(id uint64) (*model.Tag, error)
	FindAll() ([]*model.Tag, error)
	FindAllByBlog(blog *model.Blog) ([]*model.Tag, error)
	Update(tag *model.Tag) error
	Remove(tag *model.Tag) error
}
