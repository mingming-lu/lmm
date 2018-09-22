package ui

type postArticleAdaptor struct {
	Title *string  `json:"title"`
	Body  *string  `json:"body"`
	Tags  []string `json:"tags"`
}
