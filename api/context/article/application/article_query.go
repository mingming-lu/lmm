package application

import (
	"lmm/api/context/article/domain/finder"
	"lmm/api/context/article/domain/model"
	"lmm/api/utils/strings"
)

// ArticleQueryService is a query side application
type ArticleQueryService struct {
	articleFinder finder.ArticleFinder
}

// NewArticleQueryService is a constructor of ArticleQueryService
func NewArticleQueryService(articleFinder finder.ArticleFinder) *ArticleQueryService {
	return &ArticleQueryService{articleFinder: articleFinder}
}

// ListArticlesByPage is used for listing articles on article index page
func (app *ArticleQueryService) ListArticlesByPage(countStr, pageStr string) (*model.ArticleListView, error) {
	if countStr == "" {
		countStr = "5"
	}
	count, err := strings.ParseUint(countStr)
	if err != nil {
		return nil, ErrInvalidCount
	}

	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strings.ParseUint(pageStr)
	if err != nil {
		return nil, ErrInvalidPage
	}

	return app.articleFinder.ListByPage(count, page)
}
