package model

import (
	"time"

	"lmm/api/clock"
)

type ArticleID struct {
	id       int64
	authorID int64
}

func NewArticleID(id, authorID int64) *ArticleID {
	return &ArticleID{id, authorID}
}

func (id *ArticleID) ID() int64 {
	return id.id
}

func (id *ArticleID) AuthorID() int64 {
	return id.authorID
}

// Article is an aggregate root model
type Article struct {
	id           *ArticleID
	linkName     string
	content      *Content
	createdAt    time.Time
	lastModified time.Time
}

// NewArticle is a article constructor
func NewArticle(articleID *ArticleID, linkName string, content *Content, createdAt, lastModified time.Time) *Article {
	article := &Article{
		id:           articleID,
		linkName:     linkName,
		content:      content,
		createdAt:    createdAt,
		lastModified: lastModified,
	}
	return article
}

// ID returns the id of the article
func (a *Article) ID() *ArticleID {
	return a.id
}

// LinkName returns a's link name
func (a *Article) LinkName() string {
	return a.linkName
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

// LastModified time
func (a *Article) LastModified() time.Time {
	return a.lastModified
}
