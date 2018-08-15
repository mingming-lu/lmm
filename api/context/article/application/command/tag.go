package command

type NewArticleTagCommand struct {
	articleID string
	name      string
}

func (c NewArticleTagCommand) ArticleID() string {
	return c.articleID
}

func (c NewArticleTagCommand) Name() string {
	return c.name
}

type RemoveArticleTagCommand struct {
	articleID string
	name      string
}

func (c RemoveArticleTagCommand) ArticleID() string {
	return c.articleID
}

func (c RemoveArticleTagCommand) Name() string {
	return c.name
}
