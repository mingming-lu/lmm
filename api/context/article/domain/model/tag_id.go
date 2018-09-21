package model

import "lmm/api/context/base/domain/model"

// TagID is a value object to identify tag
type TagID struct {
	model.ValueObject
	articleID *ArticleID
	name      string
}

// ArticleID returns the linked article's id
func (id TagID) ArticleID() *ArticleID {
	return id.articleID
}

// Name returns the tag's name
func (id TagID) Name() string {
	return id.name
}

// Equals compares tag id with another
func (id TagID) Equals(another *TagID) bool {
	return (id.ArticleID() == another.ArticleID()) && (id.Name() == another.Name())
}
