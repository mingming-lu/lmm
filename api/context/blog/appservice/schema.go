package appservice

type Blog struct {
	ID uint64 `json:"id"`
	BlogContent
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BlogContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BlogListPage struct {
	Blog     []*Blog `json:"blog"`
	NextPage int     `json:"next_page"`
}

type Category struct {
	Name string `json:"name"`
}

type Categories struct {
	Categories []*Category `json:"categories"`
}