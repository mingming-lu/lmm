package command

import (
	account "lmm/api/context/account/domain/model"
)

type NewArticleCommand struct {
	user     *account.User
	title    string
	text     string
	tagNames []string
}

func (c NewArticleCommand) Title() string {
	return c.title
}

func (c NewArticleCommand) Text() string {
	return c.text
}

func (c NewArticleCommand) User() *account.User {
	return c.user
}

func (c NewArticleCommand) TagNames() []string {
	return c.tagNames
}

type UpdateArticleCommand struct {
	articleID    string
	articleTitle string
	articleBody  string
}

func (c UpdateArticleCommand) ArticleID() string {
	return c.articleID
}

func (c UpdateArticleCommand) ArticleTitle() string {
	return c.articleTitle
}

func (c UpdateArticleCommand) ArticleBody() string {
	return c.articleBody
}
