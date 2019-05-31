package repository

import (
	"lmm/api/pkg/transaction"

	"lmm/api/service/article/domain/model"
)

// ArticleRepository interface
type ArticleRepository interface {
	NextID(tx transaction.Transaction, authorID int64) (*model.ArticleID, error)
	Save(tx transaction.Transaction, article *model.Article) error
	Remove(tx transaction.Transaction, id *model.ArticleID) error
	FindByID(tx transaction.Transaction, id *model.ArticleID) (*model.Article, error)
}
