package application

import (
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
)

// ArticleCommandService is a command side application
type ArticleCommandService struct {
	articleService    *service.ArticleService
	articleRepository repository.ArticleRepository
	authorService     service.AuthorService
}

// NewArticleCommandService is a constructor of ArticleCommandService
func NewArticleCommandService(articleRepository repository.ArticleRepository) *ArticleCommandService {
	return &ArticleCommandService{
		articleService:    service.NewArticleService(articleRepository),
		articleRepository: articleRepository,
	}
}

// PostNewArticle is used for posting a new article
func (app *ArticleCommandService) PostNewArticle(userID uint64, title string, body string, tagNames []string) (*model.ArticleID, error) {
	author, err := app.authorService.AuthorFromUserID(userID)
	if err != nil {
		return nil, err
	}

	article, err := app.articleService.NewArticleToPost(author, title, body, tagNames)
	if err != nil {
		return nil, err
	}

	if err := app.articleRepository.Save(article); err != nil {
		return nil, err
	}

	return article.ID(), nil
}

// ModifyArticleText is used for modify the text of an article
func (app *ArticleCommandService) ModifyArticleText(rawArticleID, title, body string) error {
	article, err := app.articleWithID(rawArticleID)
	if err != nil {
		return err
	}

	newText, err := model.NewText(title, body)
	if err != nil {
		return err
	}

	if err := article.EditText(newText); err != nil {
		return err
	}

	return app.articleRepository.Save(article)
}

// AddTagToArticle is used for adding a new tag to an article
func (app *ArticleCommandService) AddTagToArticle(rawArticleID, tagName string) error {
	article, err := app.articleWithID(rawArticleID)
	if err != nil {
		return err
	}

	return app.manipulateArticleTag(article, article.AddTag, tagName)
}

// RemoveTagFromArticle is used for removing a tag from an article
func (app *ArticleCommandService) RemoveTagFromArticle(rawArticleID, tagName string) error {
	article, err := app.articleWithID(rawArticleID)
	if err != nil {
		return err
	}

	return app.manipulateArticleTag(article, article.RemoveTag, tagName)
}

func (app *ArticleCommandService) articleWithID(id string) (*model.Article, error) {
	articleID, err := model.NewArticleID(id)
	if err != nil {
		return nil, err
	}

	return app.articleRepository.FindByID(articleID)
}

func (app *ArticleCommandService) manipulateArticleTag(article *model.Article, manipulation func(*model.Tag) error, tagName string) error {
	tag, err := model.NewTag(article.ID(), tagName)
	if err != nil {
		return err
	}

	if err := manipulation(tag); err != nil {
		return err
	}

	return app.articleRepository.Save(article)
}
