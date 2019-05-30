package application

import (
	"context"
	"lmm/api/pkg/transaction"

	"github.com/pkg/errors"

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
func (app *ArticleQueryService) ListArticlesByPage(c context.Context, q query.ListArticleQuery) (*model.ArticleListView, error) {
	count, err := stringutil.ParseUint(q.Count)
	if err != nil || count < 1 {
		return nil, errors.Wrap(ErrInvalidCount, err.Error())
	}

	page, err := stringutil.ParseUint(q.Page)
	if err != nil || page < 1 {
		return nil, errors.Wrap(ErrInvalidPage, err.Error())
	}

	return nil, errors.New("not implemented")
}

// ArticleByID finds article by given id
func (app *ArticleQueryService) ArticleByID(c context.Context, rawID string) (*model.ArticleView, error) {
	panic("TODO")
	// articleID, err := model.NewArticleID(rawID)
	// if err != nil {
	// 	return nil, err
	// }
	// return app.articleFinder.FindByID(c, articleID)
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
