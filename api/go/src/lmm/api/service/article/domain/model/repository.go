package model

import "lmm/api/pkg/transaction"

// ArticleRepository interface
type ArticleRepository interface {
	NextID(tx transaction.Transaction, authorID int64) (*ArticleID, error)
	Save(tx transaction.Transaction, article *Article) error
	Remove(tx transaction.Transaction, id *ArticleID) error
	FindByID(tx transaction.Transaction, id *ArticleID) (*Article, error)
}

// ArticleViewer defines an interface to query side
type ArticleViewer interface {
	ViewArticle(tx transaction.Transaction, linkName string) (*Article, error)
	ViewArticles(tx transaction.Transaction, count, page int, filter *ArticlesFilter) (*ArticleListView, error)
	ViewAllTags(tx transaction.Transaction) ([]*TagView, error)
}

// ArticlesFilter filtering articles
type ArticlesFilter struct {
	Tag string
}
