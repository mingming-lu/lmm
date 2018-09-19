package ui

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/article/application"
	"lmm/api/context/article/domain/repository"
	"lmm/api/http"
)

// UI is the user interface to contact with network
type UI struct {
	appService *application.Service
}

// NewUI returns a new ui
func NewUI(articleRepository repository.ArticleRepository) *UI {
	appService := application.NewService(
		application.NewArticleCommandService(articleRepository),
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
		c.Header("Location", "articles/"+articleID.String()).String(http.StatusOK, "Success")
	}
}
