package command

type NewArticleCommand struct {
	title string
	text  string
}

func (c NewArticleCommand) Title() string {
	return c.title
}

func (c NewArticleCommand) Text() string {
	return c.text
}

type UpdateArticleCommand struct {
	articleID string
	title     string
	text      string
}

func (c UpdateArticleCommand) ArticleID() string {
	return c.articleID
}

func (c UpdateArticleCommand) Title() string {
	return c.title
}

func (c UpdateArticleCommand) Text() string {
	return c.text
}
