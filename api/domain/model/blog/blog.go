package blog

import (
	"lmm/api/domain/model/category"
	"lmm/api/domain/model/tag"
)

type Minimal struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ListItem struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type Blog struct {
	ID        uint64            `json:"id"`
	User      uint64            `json:"user"`
	Title     string            `json:"title"`
	Text      string            `json:"text"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Category  category.Category `json:"category"`
	Tags      []tag.Tag         `json:"tags"`
}
