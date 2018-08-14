package ui

type Image struct {
	Name string `json:"name"`
}

type ImagesPage struct {
	HasNextPage bool    `json:"has_next_page`
	Images      []Image `json:"images"`
}
