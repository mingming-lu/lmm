package model

import (
	"time"

	"lmm/api/clock"
	"lmm/api/service/base/model"
)

type ArticleID int64

// Article is an aggregate root model
type Article struct {
	model.Entity
	id           ArticleID
	author       *Author
	content      *Content
	createdAt    time.Time
	lastModified time.Time
}

// NewArticle is a article constructor
func NewArticle(articleID ArticleID, author *Author, content *Content, createdAt, lastModified time.Time) *Article {
	article := &Article{
		id:           articleID,
		author:       author,
		content:      content,
		createdAt:    createdAt,
		lastModified: lastModified,
	}
	return article
}

// ID returns the id of the article
func (a *Article) ID() ArticleID {
	return a.id
}

// Author returns the author of the article
func (a *Article) Author() *Author {
	return a.author
}

// Content returns article's content
func (a *Article) Content() *Content {
	return a.content
}

// EditContent changes article's content and update lastModified if text differs
func (a *Article) EditContent(content *Content) {
	if a.content.Text() != content.Text() {
		a.lastModified = clock.Now()
	}
	a.content = content
}

// CreatedAt time
func (a *Article) CreatedAt() time.Time {
	return a.createdAt
}

// LastModified time
func (a *Article) LastModified() time.Time {
	return a.lastModified
}
