package blog

type Minimal struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Blog struct {
	ID        int64  `json:"id"`
	User      int64  `json:"user"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
