package finder

import "lmm/api/context/article/domain/model"

// ArticleFinder defines an interface to query side
type ArticleFinder interface {
	ListByPage(count, page uint) (*model.ArticleListView, error)
	FindByID(id *model.ArticleID) (*model.ArticleView, error)
}
