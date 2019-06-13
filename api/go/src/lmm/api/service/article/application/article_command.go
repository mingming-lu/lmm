package application

import (
	"context"

	"lmm/api/clock"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/application/command"
	"lmm/api/service/article/domain/model"
	"lmm/api/util/stringutil"

	"github.com/pkg/errors"
)

// ArticleCommandService is a command side application
type ArticleCommandService struct {
	articleRepository  model.ArticleRepository
	transactionManager transaction.Manager
}

// NewArticleCommandService is a constructor of ArticleCommandService
func NewArticleCommandService(articleRepository model.ArticleRepository, transactionManager transaction.Manager) *ArticleCommandService {
	return &ArticleCommandService{
		articleRepository:  articleRepository,
		transactionManager: transactionManager,
	}
}

// PostNewArticle is used for posting a new article
func (app *ArticleCommandService) PostNewArticle(c context.Context, cmd command.PostArticle) (id *model.ArticleID, err error) {
	content, err := model.NewContent(cmd.Title, cmd.Body, cmd.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "invalid article content")
	}

	err = app.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		now := clock.Now()

		id, err = app.articleRepository.NextID(tx, cmd.AuthorID)
		if err != nil {
			return err
		}

		article := model.NewArticle(id, stringutil.Int64ToStr(id.ID()), content, now, now)

		return app.articleRepository.Save(tx, article)
	}, nil)

	return
}

// EditArticle command
func (app *ArticleCommandService) EditArticle(c context.Context, cmd command.EditArticle) error {
	articleID := model.NewArticleID(cmd.ArticleID, cmd.UserID)

	content, err := model.NewContent(cmd.Title, cmd.Body, cmd.Tags)
	if err != nil {
		return errors.Wrap(err, "invalid article content")
	}

	return app.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		article, err := app.articleRepository.FindByID(tx, articleID)
		if err != nil {
			return errors.Wrap(err, "article not found")
		}

		if err := article.ChangeLinkName(cmd.LinkName); err != nil {
			return errors.Wrap(err, "invalid article link name")
		}
		article.EditContent(content)

		return app.articleRepository.Save(tx, article)
	}, nil)
}
