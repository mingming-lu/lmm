package ui

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/article/application"
	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
	"lmm/api/http"
)

// UI is the user interface to contact with network
type UI struct {
	appService *application.Service
}

// NewUI returns a new ui
func NewUI(articleRepository repository.ArticleRepository, authorService service.AuthorService) *UI {
	appService := application.NewService(
		application.NewArticleCommandService(articleRepository, authorService),
		nil,
	)
	return &UI{appService: appService}
}

// PostArticle handles POST /1/articles
func (ui *UI) PostArticle(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	article := PostingArticleAdaptor{}
	if err := c.Request.ScanBody(&article); err != nil {
		http.BadRequest(c)
		return
	}

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(
		user.ID(),
		article.Title,
		article.Body,
		article.Tags,
	)
	switch err {
	case nil:
		c.Header("Location", "articles/"+articleID.String()).String(http.StatusCreated, "Success")
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrNoSuchUser:
		http.Unauthorized(c)
	default:
		panic(err)
	}
}

// EditArticleText handles PUT /1/article/:articleID
func (ui *UI) EditArticleText(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	article := EditArticleAdaptor{}
	if err := c.Request.ScanBody(&article); err != nil {
		http.BadRequest(c)
		return
	}

	err := ui.appService.ArticleCommandService().EditArticleText(
		user.ID(),
		c.Request.Path.Params("articleID"),
		article.Title,
		article.Body,
	)
	switch err {
	case nil:
		http.NoContent(c)
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrNoSuchArticle:
		c.String(http.StatusNotFound, err.Error())
	case domain.ErrNoSuchUser:
		http.Unauthorized(c)
	case domain.ErrNotArticleAuthor:
		c.String(http.StatusForbidden, err.Error())
	default:
		panic(err)
	}
}
