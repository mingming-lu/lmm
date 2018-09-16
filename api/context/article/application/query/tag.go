package query

import "lmm/api/context/article/domain/model"

type TagFinder interface {
	FindAllTags() ([]*model.Tag, error)
}
