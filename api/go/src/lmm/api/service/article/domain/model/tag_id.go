package model

import "lmm/api/model"

// TagID is a value object to identify tag
type TagID struct {
	model.ValueObject
	articleID *ArticleID
	order     uint
}

// ArticleID returns the linked article's id
func (id *TagID) ArticleID() *ArticleID {
	return id.articleID
}

// Order returns the tag's name
func (id *TagID) Order() uint {
	return id.order
}

// Equals compares tag id with another
func (id *TagID) Equals(another *TagID) bool {
	return (id.ArticleID() == another.ArticleID()) && (id.Order() == another.Order())
}
