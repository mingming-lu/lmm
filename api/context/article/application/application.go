package application

import (
	"lmm/api/context/article/application/command"
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
)

type AppService struct {
	articlePostingService service.ArticlePostingService
	articleRepo           repository.ArticleRepository
}

func (app *AppService) NewArticle(cmd command.NewArticleCommand) (string, error) {
	article, err := app.articlePostingService.PostingArticle(
		cmd.User(),
		cmd.Title(),
		cmd.Text(),
		cmd.TagNames(),
	)
	if err != nil {
		return "", err
	}

	if err := app.articleRepo.Save(article); err != nil {
		return "", err
	}

	return article.ID().String(), nil
}

func (app *AppService) UpdateArticle(cmd command.UpdateArticleCommand) error {
	articleID, err := model.NewArticleID(cmd.ArticleID())
	if err != nil {
		return err
	}

	article, err := model.NewArticle(articleID, cmd.Title(), cmd.Text())
	if err != nil {
		return err
	}

	return app.articleRepo.Update(article)
}

func (app *AppService) NewArticleTag(cmd command.NewArticleTagCommand) (string, error) {
	articleID, err := model.NewArticleID(cmd.ArticleID())
	if err != nil {
		return "", err
	}

	tag, err := model.NewTag(articleID, cmd.Name())
	if err != nil {
		return "", err
	}

	if err := app.tagRepo.Save(tag); err != nil {
		return "", err
	}

	return tag.ID().Name(), nil
}

func (app *AppService) RemoveArticleTag(cmd command.RemoveArticleTagCommand) error {
	articleID, err := model.NewArticleID(cmd.ArticleID())
	if err != nil {
		return err
	}

	tag, err := model.NewTag(articleID, cmd.Name())
	if err != nil {
		return err
	}

	return app.tagRepo.Remove(tag.ID())
}
