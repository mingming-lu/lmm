package application

import (
	"lmm/api/context/article/application/command"
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
)

// ArticleApplicationService is a command side application
type ArticleApplicationService struct {
	articlePostingService *service.ArticlePostingService
	articleRepository     repository.ArticleRepository
}

// NewArticleApplicationService is a constructor of ArticleApplicationService
func NewArticleApplicationService(articleRepository repository.ArticleRepository) *ArticleApplicationService {
	return &ArticleApplicationService{
		articlePostingService: service.NewArticlePostingService(articleRepository),
		articleRepository:     articleRepository,
	}
}

// PostingArticle is used for posting a new article
func (app *ArticleApplicationService) PostingArticle(cmd command.PostingArticleCommand) (string, error) {
	article, err := app.articlePostingService.PostingArticle(
		cmd.User(),
		cmd.Title(),
		cmd.Body(),
		cmd.TagNames(),
	)
	if err != nil {
		return "", err
	}

	if err := app.articleRepository.Save(article); err != nil {
		return "", err
	}

	return article.ID().String(), nil
}

// ModifyArticleText is used for modify the text of an article
func (app *ArticleApplicationService) ModifyArticleText(cmd command.ModifyArticleCommand) error {
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

	return app.articleRepository.Save(article)
}

// NewArticleTag is used for adding a new tag to an article
func (app *ArticleApplicationService) NewArticleTag(cmd command.NewArticleTagCommand) error {
	article, err := app.articleWithID(cmd.ArticleID())
	if err != nil {
		return err
	}

	return app.manipulateArticleTag(cmd.TagName(), article, article.AddTag)
}

// RemoveArticleTag is used for removing a tag from an article
func (app *ArticleApplicationService) RemoveArticleTag(cmd command.RemoveArticleTagCommand) error {
	article, err := app.articleWithID(cmd.ArticleID())
	if err != nil {
		return err
	}

	return app.manipulateArticleTag(cmd.TagName(), article, article.RemoveTag)
}

func (app *ArticleApplicationService) articleWithID(id string) (*model.Article, error) {
	articleID, err := model.NewArticleID(id)
	if err != nil {
		return nil, err
	}

	return app.articleRepository.FindByID(articleID)
}

func (app *ArticleApplicationService) manipulateArticleTag(
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

	return app.articleRepository.Save(article)
}
