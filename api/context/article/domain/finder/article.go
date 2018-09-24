package finder

import "lmm/api/context/article/domain/model"

// ArticleFinder defines an interface to query side
type ArticleFinder interface {
	FindByID(id *model.ArticleID) (*model.ArticleView, error)
	ListAllTags() (model.TagListView, error)
	ListByPage(count, page uint) (*model.ArticleListView, error)
}
