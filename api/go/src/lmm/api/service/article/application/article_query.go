package application

import (
	"context"
	"lmm/api/service/article/application/query"

	"lmm/api/service/article/domain/finder"
	"lmm/api/service/article/domain/model"
	"lmm/api/util/stringutil"
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
func (app *ArticleQueryService) ListArticlesByPage(c context.Context, q query.ListArticleQuery) (*model.ArticleListView, error) {
	var (
		count uint = 5
		page  uint = 1
	)
	if q.Count != "" && q.Count != "5" {
		count, err := stringutil.ParseUint(q.Count)
		if err != nil || count < 1 {
			return nil, ErrInvalidCount
		}
	}

	if q.Page != "" && q.Page != "1" {
		page, err := stringutil.ParseUint(q.Page)
		if err != nil || page < 1 {
			return nil, ErrInvalidPage
		}
	}

	return app.articleFinder.ListByPage(c, count, page, &finder.ArticleFilter{
		Tags: q.Tags,
	})
}

// ArticleByID finds article by given id
func (app *ArticleQueryService) ArticleByID(c context.Context, rawID string) (*model.ArticleView, error) {
	articleID, err := model.NewArticleID(rawID)
	if err != nil {
		return nil, err
	}
	return app.articleFinder.FindByID(c, articleID)
}

// AllArticleTags gets all article tags
func (app *ArticleQueryService) AllArticleTags(c context.Context) (model.TagListView, error) {
	return app.articleFinder.ListAllTags(c)
}
