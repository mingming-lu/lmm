package finder

import (
	"context"

	"lmm/api/service/article/domain/model"
)

// ArticleFinder defines an interface to query side
type ArticleFinder interface {
	FindByID(c context.Context, id *model.ArticleID) (*model.ArticleView, error)
	ListAllTags(c context.Context) (model.TagListView, error)
	ListByPage(c context.Context, count, page uint) (*model.ArticleListView, error)
}
