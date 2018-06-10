package ui

type Blog struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BlogListResponse struct {
	Blog        []BlogResponse `json:"blog"`
	Page        int            `json:"page"`
	HasNextPage bool           `json:"has_next_page"`
}

type BlogResponse struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
