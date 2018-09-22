package ui

import (
	"errors"

	account "lmm/api/context/account/domain/model"
	"lmm/api/context/article/application"
	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/repository"
	"lmm/api/context/article/domain/service"
	"lmm/api/http"
)

var (
	errTitleRequired = errors.New("title required")
	errBodyRequired  = errors.New("body requried")
	errTagsRequired  = errors.New("tags requried")
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

	article := postArticleAdaptor{}
	if err := c.Request.ScanBody(&article); err != nil {
		http.BadRequest(c)
		return
	}

	if err := ui.validatePostArticleAdaptor(&article); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(
		user.ID(),
		*article.Title,
		*article.Body,
		article.Tags,
	)
	switch err {
	case nil:
		c.Header("Location", "/v1/articles/"+articleID.String()).String(http.StatusCreated, "Success")
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

	article := postArticleAdaptor{}
	if err := c.Request.ScanBody(&article); err != nil {
		http.BadRequest(c)
		return
	}

	err := ui.appService.ArticleCommandService().EditArticle(
		user.ID(),
		c.Request.Path.Params("articleID"),
		*article.Title,
		*article.Body,
		article.Tags,
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

func (ui *UI) validatePostArticleAdaptor(adaptor *postArticleAdaptor) error {
	if adaptor.Title == nil {
		return errTitleRequired
	}
	if adaptor.Body == nil {
		return errBodyRequired
	}
	if adaptor.Tags == nil {
		return errTagsRequired
	}
	return nil
}
