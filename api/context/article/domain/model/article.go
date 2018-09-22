package model

import (
	"errors"

	"lmm/api/context/base/domain/model"
)

var (
	ErrArticleTextNoChange = errors.New("article text no change")
	ErrTagAlreadyAdded     = errors.New("the tag has already been added")
	ErrTagNotExists        = errors.New("tag not exists")
)

// Article is an aggregate root model
type Article struct {
	model.Entity
	id      *ArticleID
	author  *Author
	content *Content
}

// NewArticle is a article constructor
func NewArticle(articleID *ArticleID, author *Author, content *Content) *Article {
	return &Article{
		id:      articleID,
		author:  author,
		content: content,
	}
}

// ID returns the id of the article
func (a *Article) ID() *ArticleID {
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

// EditContent changes article's content
func (a *Article) EditContent(content *Content) {
	a.content = content
}
