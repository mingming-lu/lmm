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
	Page      int               `json:"page"`
	PerPage   int               `json:"perPage"`
	Tag       string            `json:"tag,omitempty"`
	Total     int               `json:"total"`
	PrevPage  string            `json:"prevPage,omitempty"`
	NextPage  string            `json:"nextPage,omitempty"`
	FirstPage string            `json:"firstPage,omitempty"`
	LastPage  string            `json:"lastPage,omitempty"`
}

type articleListItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	PostAt int64  `json:"post_at,string"`
}

type articleViewResponse struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Body         string           `json:"body"`
	PostAt       int64            `json:"post_at,string"`
	LastEditedAt int64            `json:"last_edited_at,string"`
	Tags         []articleViewTag `json:"tags"`
}

type articleViewTag struct {
	Name string `json:"name"`
}

type articleTagListView = []*articleTagListItemView

type articleTagListItemView struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
