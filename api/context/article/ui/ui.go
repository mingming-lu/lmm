package ui

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/article/application"
	"lmm/api/context/article/application/command"
	"lmm/api/http"
)

type UI struct {
	app *application.QueryAppService
}

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

	articleID, err := ui.app.PostingArticle(command.NewPostingArticleCommand(
		user,
		article.Title,
		article.Body,
		article.Tags,
	))
	switch err {
	case nil:
		c.Header("Location", "articles/"+articleID).String(http.StatusOK, "Success")
	}
}

func (ui *UI) ModifyArticle(c *http.Context) {
	user, ok := c.Values().Get("user").(*account.User)
	if !ok {
		http.Unauthorized(c)
		return
	}

	article := ModifyArticleAdaptor{}
	if err := c.Request.ScanBody(&article); err != nil {
		http.BadRequest(c)
		return
	}

	err := ui.app.ModifyArticleText(command.NewModifyArticleCommand(
		user,
		c.Request.Path.Params("article"),
		article.Title,
		article.Body,
	))
	switch err {
	case nil:
		c.String(http.StatusOK, "Success")
	}
}
