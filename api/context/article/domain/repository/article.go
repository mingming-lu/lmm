package repository

import "lmm/api/context/article/domain/model"

type ArticleRepository interface {
	NextID() string
	Save(*model.Article) error
	Remove(*model.Article) error
	FindByID(model.ArticleID) (*model.Article, error)
}
