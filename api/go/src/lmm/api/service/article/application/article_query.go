package application

import (
	"context"

	"github.com/pkg/errors"

	"lmm/api/pkg/transaction"
	"lmm/api/service/article/application/query"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/viewer"
	"lmm/api/util/stringutil"
)

// ArticleQueryService is a query side application
type ArticleQueryService struct {
	viewer    viewer.ArticleViewer
	txManager transaction.Manager
}

// NewArticleQueryService is a constructor of ArticleQueryService
func NewArticleQueryService(viewer viewer.ArticleViewer, txManager transaction.Manager) *ArticleQueryService {
	return &ArticleQueryService{viewer: viewer, txManager: txManager}
}

// ListArticlesByPage is used for listing articles on article index page
func (app *ArticleQueryService) ListArticlesByPage(c context.Context, q query.ListArticleQuery) (articles *model.ArticleListView, err error) {
	count, err := stringutil.ParseInt(q.Count)
	if err != nil || count < 1 {
		return nil, errors.Wrap(ErrInvalidCount, err.Error())
	}

	page, err := stringutil.ParseInt(q.Page)
	if err != nil || page < 1 {
		return nil, errors.Wrap(ErrInvalidPage, err.Error())
	}

	err = app.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		articles, err = app.viewer.ViewArticles(tx, count, page, nil)

		return err
	}, &transaction.Option{ReadOnly: true})

	return
}

func (app *ArticleQueryService) ArticleByLinkName(c context.Context, linkName string) (article *model.Article, err error) {
	err = app.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		article, err = app.viewer.ViewArticle(tx, linkName)
		return err
	}, &transaction.Option{ReadOnly: true})

	return
}

// AllArticleTags gets all article tags
func (app *ArticleQueryService) AllArticleTags(c context.Context) (tags []*model.TagView, err error) {
	err = app.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		tags, err = app.viewer.ViewAllTags(tx)
		if err != nil {
			return errors.Wrap(err, "unexpect error when fetch all tags")
		}
		return err
	}, &transaction.Option{ReadOnly: true})

	return
}
