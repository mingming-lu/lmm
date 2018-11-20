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

type articleListAdapterV2 struct {
	Articles  []articleListItem `json:"articles"`
	Page      uint              `json:"page"`
	PerPage   uint              `json:"perPage"`
	Total     uint              `json:"total"`
	PrevPage  *string           `json:"prevPage"`
	NextPage  *string           `json:"nextPage"`
	FirstPage *string           `json:"firstPage"`
	LastPage  *string           `json:"lastPage"`
}

type articleListItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	PostAt string `json:"post_at"`
}

type articleViewResponse struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Body         string           `json:"body"`
	PostAt       string           `json:"post_at"`
	LastEditedAt string           `json:"last_edited_at"`
	Tags         []articleViewTag `json:"tags"`
}

type articleViewTag struct {
	Name string `json:"name"`
}

type articleTagListView = []articleTagListItemView

type articleTagListItemView struct {
	Name string `json:"name"`
}
