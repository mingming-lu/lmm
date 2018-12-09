package application

import (
	"context"

	"github.com/pkg/errors"

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
	count, err := stringutil.ParseUint(q.Count)
	if err != nil || count < 1 {
		return nil, errors.Wrap(ErrInvalidCount, err.Error())
	}

	page, err := stringutil.ParseUint(q.Page)
	if err != nil || page < 1 {
		return nil, errors.Wrap(ErrInvalidPage, err.Error())
	}

	return app.articleFinder.ListByPage(c, count, page, finder.ArticleFilter{
		Tag: q.Tag,
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
