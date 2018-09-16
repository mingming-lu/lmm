package query

import "lmm/api/context/article/domain/model"

type ArticleFinder interface {
	FindArticlesByPage(count, page uint) ([]*model.Article, bool, error)
}
