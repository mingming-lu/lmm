package application

import (
	"context"
	"lmm/api/service/article/application/command"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/service"
)

// ArticleCommandService is a command side application
type ArticleCommandService struct {
	articleService    *service.ArticleService
	articleRepository repository.ArticleRepository
	authorService     service.AuthorService
}

// NewArticleCommandService is a constructor of ArticleCommandService
func NewArticleCommandService(articleRepository repository.ArticleRepository, authorService service.AuthorService) *ArticleCommandService {
	return &ArticleCommandService{
		articleService:    service.NewArticleService(articleRepository),
		articleRepository: articleRepository,
		authorService:     authorService,
	}
}

// PostNewArticle is used for posting a new article
func (app *ArticleCommandService) PostNewArticle(c context.Context, userName, title string, body string, tagNames []string) (*model.ArticleID, error) {
	author, err := app.authorService.AuthorFromUserName(c, userName)
	if err != nil {
		return nil, err
	}

	article, err := app.articleService.NewArticleToPost(c, author, title, body, tagNames)
	if err != nil {
		return nil, err
	}

	if err := app.articleRepository.Save(c, article); err != nil {
		return nil, err
	}

	return article.ID(), nil
}

// EditArticle is used for edit the article content
func (app *ArticleCommandService) EditArticle(c context.Context, cmd *command.EditArticle) error {
	author, err := app.authorService.AuthorFromUserName(c, cmd.UserName)
	if err != nil {
		return err
	}

	article, err := app.articleWithID(c, cmd.TargetArticleID)
	if err != nil {
		return err
	}

	if cmd.AliasArticleID != "" {
		if err := article.ID().SetAlias(cmd.AliasArticleID); err != nil {
			return err
		}
	}

	if article.Author().Name() != author.Name() {
		return domain.ErrNotArticleAuthor
	}

	newText, err := model.NewText(cmd.Title, cmd.Body)
	if err != nil {
		return err
	}

	newTags, err := app.tagsFromNames(cmd.TagNames, article.ID())
	if err != nil {
		return err
	}

	content, err := model.NewContent(newText, newTags)

	article.EditContent(content)

	return app.articleRepository.Save(c, article)
}

func (app *ArticleCommandService) articleWithID(c context.Context, id string) (*model.Article, error) {
	articleID, err := model.NewArticleID(id)
	if err != nil {
		return nil, err
	}

	return app.articleRepository.FindByID(c, articleID)
}

func (app *ArticleCommandService) tagsFromNames(tagNames []string, articleID *model.ArticleID) ([]*model.Tag, error) {
	tags := make([]*model.Tag, len(tagNames), len(tagNames))

	for i, name := range tagNames {
		tag, err := model.NewTag(articleID, uint(i+1), name)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	return tags, nil
}
