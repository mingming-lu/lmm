package command

// EditArticle command
type EditArticle struct {
	UserName        string
	TargetArticleID string
	AliasArticleID  string
	Title           string
	Body            string
	TagNames        []string
}
