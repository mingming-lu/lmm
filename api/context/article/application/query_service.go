package application

import (
	"lmm/api/context/article/application/query"
	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/model"
	"lmm/api/utils/strings"
)

type QueryAppService struct {
	articleFinder query.ArticleFinder
	tagFinder     query.TagFinder
}

func NewQueryAppService(
	articleFinder query.ArticleFinder,
	tagFinder query.TagFinder,
) *QueryAppService {
	return &QueryAppService{articleFinder: articleFinder, tagFinder: tagFinder}
}

func (app *QueryAppService) ListArticlesByPage(countStr, pageStr string) ([]*model.Article, bool, error) {
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

func (app *QueryAppService) ListAllTags() ([]*model.Tag, error) {
	return app.tagFinder.FindAllTags()
}
