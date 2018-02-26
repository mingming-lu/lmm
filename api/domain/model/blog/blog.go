package blog

type Minimal struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ListItem struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type Blog struct {
	ID        int64  `json:"id"`
	User      int64  `json:"user"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
