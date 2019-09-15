package model

import (
	"time"

	"lmm/api/clock"
)

type ArticleID string

func NewArticleID(s string) *ArticleID {
	id := ArticleID(s)
	return &id
}

func (id *ArticleID) String() string {
	return string(*id)
}

// Article is an aggregate root model
type Article struct {
	id           *ArticleID
	author       *Author
	linkName     string
	content      *Content
	createdAt    time.Time
	publishedAt  time.Time
	lastModified time.Time
}

// NewArticle is a article constructor
func NewArticle(articleID *ArticleID, author *Author, content *Content, createdAt, publishedAt, lastModified time.Time) *Article {
	article := &Article{
		id:           articleID,
		author:       author,
		content:      content,
		createdAt:    createdAt,
		publishedAt:  publishedAt,
		lastModified: lastModified,
	}
	return article
}

// ID returns the id of the article
func (a *Article) ID() *ArticleID {
	return a.id
}

func (a *Article) Author() *Author {
	return a.author
}

// ChangeLinkName changed a's LinkName to newLinkName
// TODO: validate newLinkName
func (a *Article) ChangeLinkName(newLinkName string) error {
	a.linkName = newLinkName

	return nil
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

// Published returns true if a is published
// and the PublishedAt should return a non-zero value
func (a *Article) Published() bool {
	return !a.publishedAt.IsZero()
}

// PublishedAt time
func (a *Article) PublishedAt() time.Time {
	return a.publishedAt
}

// LastModified time
func (a *Article) LastModified() time.Time {
	return a.lastModified
}
