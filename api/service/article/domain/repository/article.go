package repository

import (
	"context"

	"lmm/api/context/article/domain/model"
)

// ArticleRepository interface
type ArticleRepository interface {
	NextID(context.Context) string
	Save(context.Context, *model.Article) error
	Remove(context.Context, *model.Article) error
	FindByID(context.Context, *model.ArticleID) (*model.Article, error)
}
