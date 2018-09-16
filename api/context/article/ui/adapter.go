package ui

type PostingArticleAdaptor struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type ModifyArticleAdaptor struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
