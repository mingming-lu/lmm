package category

type Minimal struct {
	Name string `json:"name"`
}

type Category struct {
	ID   int64  `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
}
