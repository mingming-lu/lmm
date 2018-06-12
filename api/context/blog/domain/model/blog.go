package model

import (
	"lmm/api/domain/model"
	"time"
)

type Blog struct {
	model.Entity
	id        uint64
	user      uint64
	title     string
	text      string
	createdAt time.Time
	updatedAt time.Time
}

func NewBlog(id, userID uint64, title, text string, createdAt, updatedAt time.Time) *Blog {
	return &Blog{
		id:        id,
		user:      userID,
		title:     title,
		text:      text,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (b *Blog) update() {
	b.updatedAt = time.Now()
}

func (b *Blog) UpdateTitle(title string) {
	if title == "" {
		return
	}
	if b.title == title {
		return
	}
	b.title = title
	b.update()
}

func (b *Blog) UpdateText(text string) {
	if b.text == text {
		return
	}
	b.text = text
	b.update()
}

func (b *Blog) ID() uint64 {
	return b.id
}

func (b *Blog) UserID() uint64 {
	return b.user
}

func (b *Blog) Title() string {
	return b.title
}

func (b *Blog) Text() string {
	return b.text
}

func (b *Blog) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Blog) UpdatedAt() time.Time {
	return b.updatedAt
}
