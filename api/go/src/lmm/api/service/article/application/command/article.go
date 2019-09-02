package command

// PostArticle Command
type PostArticle struct {
	AuthorID int64
	Title    string
	Body     string
	Tags     []string
}

// EditArticle command
type EditArticle struct {
	UserID    int64
	ArticleID string
	LinkName  string
	Title     string
	Body      string
	Tags      []string
}

// PublishArticle command
type PublishArticle struct {
	UserID    int64
	ArticleID string
}
