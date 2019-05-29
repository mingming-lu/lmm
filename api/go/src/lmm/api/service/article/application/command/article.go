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
	UserName        string
	TargetArticleID string
	AliasArticleID  string
	Title           string
	Body            string
	TagNames        []string
}
