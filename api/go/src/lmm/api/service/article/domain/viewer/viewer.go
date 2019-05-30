package viewer

import (
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/domain/model"
)

// ArticleViewer defines an interface to query side
type ArticleViewer interface {
	ViewArticle(tx transaction.Transaction, id model.ArticleID) (*model.ArticleView, error)
	ViewArticles(tx transaction.Transaction, count, page uint, filter *ArticlesFilter) (*model.ArticleListView, error)
	ViewAllTags(tx transaction.Transaction) ([]*model.TagView, error)
}

// ArticlesFilter filtering articles
type ArticlesFilter struct {
	Tag string
}
