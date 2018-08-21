package application

import (
	"lmm/api/context/article/application/command"
	"lmm/api/context/article/application/query"
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
)

type CommandAppService struct {
	articlePostingService *service.ArticlePostingService
	articleRepo           repository.ArticleRepository
}

func NewCommandAppService(
	articleFinder query.ArticleFinder,
	articleRepo repository.ArticleRepository,
) *CommandAppService {
	return &CommandAppService{
		articlePostingService: service.NewArticlePostingService(articleRepo),
		articleRepo:           articleRepo,
	}
}

func (app *CommandAppService) PostingArticle(cmd command.PostingArticleCommand) (string, error) {
	article, err := app.articlePostingService.PostingArticle(
		cmd.User(),
		cmd.Title(),
		cmd.Body(),
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

func (app *CommandAppService) ModifyArticleText(cmd command.ModifyArticleCommand) error {
	article, err := app.articleWithID(cmd.ArticleID())
	if err != nil {
		return err
	}

	newArticleText, err := model.NewArticleText(cmd.ArticleTitle(), cmd.ArticleBody())
	if err != nil {
		return err
	}

	if err := article.ModifyText(newArticleText); err != nil {
		return err
	}

	return app.articleRepo.Save(article)
}

func (app *CommandAppService) NewArticleTag(cmd command.NewArticleTagCommand) error {
	article, err := app.articleWithID(cmd.ArticleID())
	if err != nil {
		return err
	}

	return app.manipulateArticleTag(cmd.TagName(), article, article.AddTag)
}

func (app *CommandAppService) RemoveArticleTag(cmd command.RemoveArticleTagCommand) error {
	article, err := app.articleWithID(cmd.ArticleID())
	if err != nil {
		return err
	}

	return app.manipulateArticleTag(cmd.TagName(), article, article.RemoveTag)
}

func (app *CommandAppService) articleWithID(id string) (*model.Article, error) {
	articleID, err := model.NewArticleID(id)
	if err != nil {
		return nil, err
	}

	return app.articleRepo.FindByID(articleID)
}

func (app *CommandAppService) manipulateArticleTag(
	tagName string,
	article *model.Article,
	manipulation func(*model.Tag) error,
) error {
	tag, err := model.NewTag(article.ID(), tagName)
	if err != nil {
		return err
	}

	if err := manipulation(tag); err != nil {
		return err
	}

	return app.articleRepo.Save(article)
}
