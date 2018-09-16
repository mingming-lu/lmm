package command

import (
	account "lmm/api/context/account/domain/model"
)

type PostingArticleCommand struct {
	user     *account.User
	title    string
	body     string
	tagNames []string
}

func NewPostingArticleCommand(user *account.User, title, body string, tagNames []string) PostingArticleCommand {
	return PostingArticleCommand{
		user:     user,
		title:    title,
		body:     body,
		tagNames: tagNames,
	}
}

func (c PostingArticleCommand) Title() string {
	return c.title
}

func (c PostingArticleCommand) Body() string {
	return c.body
}

func (c PostingArticleCommand) User() *account.User {
	return c.user
}

func (c PostingArticleCommand) TagNames() []string {
	return c.tagNames
}

type ModifyArticleCommand struct {
	user         *account.User
	articleID    string
	articleTitle string
	articleBody  string
}

func NewModifyArticleCommand(user *account.User, id, title, body string) ModifyArticleCommand {
	return ModifyArticleCommand{
		user:         user,
		articleID:    id,
		articleTitle: title,
		articleBody:  body,
	}
}

func (c ModifyArticleCommand) User() *account.User {
	return c.user
}

func (c ModifyArticleCommand) ArticleID() string {
	return c.articleID
}

func (c ModifyArticleCommand) ArticleTitle() string {
	return c.articleTitle
}

func (c ModifyArticleCommand) ArticleBody() string {
	return c.articleBody
}
