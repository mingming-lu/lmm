package ui

type PostingArticleAdaptor struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type EditArticleAdaptor struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
