package application

import (
	"context"

	"lmm/api/clock"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/application/command"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
)

// ArticleCommandService is a command side application
type ArticleCommandService struct {
	articleRepository  repository.ArticleRepository
	transactionManager transaction.Manager
}

// NewArticleCommandService is a constructor of ArticleCommandService
func NewArticleCommandService(articleRepository repository.ArticleRepository, transactionManager transaction.Manager) *ArticleCommandService {
	return &ArticleCommandService{
		articleRepository:  articleRepository,
		transactionManager: transactionManager,
	}
}

// PostNewArticle is used for posting a new article
func (app *ArticleCommandService) PostNewArticle(c context.Context, cmd command.PostArticle) (id model.ArticleID, err error) {
	author := model.NewAuthor(cmd.AuthorID)

	text, err := model.NewText(cmd.Title, cmd.Body)
	if err != nil {
		return -1, err
	}

	content := model.NewContent(text, nil)

	err = app.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		now := clock.Now()

		id, err = app.articleRepository.NextID(tx, cmd.AuthorID)
		if err != nil {
			return err
		}

		article := model.NewArticle(id, author, content, now, now)

		return app.articleRepository.Save(tx, article)
	}, nil)

	return
}
