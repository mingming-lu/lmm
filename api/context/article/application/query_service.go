package application

import (
	"lmm/api/context/article/application/query"
	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/model"
	"lmm/api/utils/strings"
)

// ArticleQueryService is a query side application
type ArticleQueryService struct {
	articleFinder query.ArticleFinder
	tagFinder     query.TagFinder
}

// NewArticleQueryService is a constructor of ArticleQueryService
func NewArticleQueryService(articleFinder query.ArticleFinder, tagFinder query.TagFinder) *ArticleQueryService {
	return &ArticleQueryService{articleFinder: articleFinder, tagFinder: tagFinder}
}

// ListArticlesByPage is used for listing articles on article index page
func (app *ArticleQueryService) ListArticlesByPage(countStr, pageStr string) ([]*model.Article, bool, error) {
	if countStr == "" {
		countStr = "10"
	}
	count, err := strings.ParseUint(countStr)
	if err != nil {
		return nil, false, domain.ErrInvalidCount
	}

	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strings.ParseUint(pageStr)
	if err != nil {
		return nil, false, domain.ErrInvalidPage
	}

	return app.articleFinder.FindArticlesByPage(count, page)
}

// ListAllTags is used for listing all tags on article index page
func (app *ArticleQueryService) ListAllTags() ([]*model.Tag, error) {
	return app.tagFinder.FindAllTags()
}
