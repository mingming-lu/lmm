package project

type Project struct {
	ID          int64   `json:"id"`
	User        int64   `json:"user"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	FromDate    *string `json:"from_date"`
	ToDate      *string `json:"to_date"`
}
