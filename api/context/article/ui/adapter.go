package ui

type postArticleAdapter struct {
	Title *string  `json:"title"`
	Body  *string  `json:"body"`
	Tags  []string `json:"tags"`
}

type articleListAdapter struct {
	Articles    []articleListItem `json:"articles"`
	HasNextPage bool              `json:"has_next_page"`
}

type articleListItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	PostAt string `json:"post_at"`
}
