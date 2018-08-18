package repository

import "lmm/api/context/article/domain/model"

type ArticleRepository interface {
	NextID() string
	Save(*model.Article) error
	Update(*model.Article) error
	Find(model.ArticleID) (*model.Article, error)
}
