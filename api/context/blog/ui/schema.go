package ui

type Blog struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BlogResponse struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
