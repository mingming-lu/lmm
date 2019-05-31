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
	ArticleID int64
	Title     string
	Body      string
	Tags      []string
}
