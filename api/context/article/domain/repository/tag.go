package repository

import "lmm/api/context/article/domain/model"

type TagRepository interface {
	Save(*model.Tag) error
	Remove(model.TagID) error
}
