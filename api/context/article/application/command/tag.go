package command

type NewArticleTagCommand struct {
	articleID string
	tagName   string
}

func (c NewArticleTagCommand) ArticleID() string {
	return c.articleID
}

func (c NewArticleTagCommand) TagName() string {
	return c.tagName
}

type RemoveArticleTagCommand struct {
	articleID string
	tagName   string
}

func (c RemoveArticleTagCommand) ArticleID() string {
	return c.articleID
}

func (c RemoveArticleTagCommand) TagName() string {
	return c.tagName
}
