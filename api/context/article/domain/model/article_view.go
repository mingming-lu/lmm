package model

import (
	"time"

	"lmm/api/context/base/domain/model"
)

// ArticleView model defines what should be shown on article page
type ArticleView struct {
	model.ValueObject
	id           *ArticleID
	author       *Author
	content      *Content
	postAt       time.Time
	lastEditedAt time.Time
}

// ID returns article's id
func (v *ArticleView) ID() *ArticleID {
	return v.id
}

// Author returns article's author
func (v *ArticleView) Author() *Author {
	return v.author
}

// Content returns article's content
func (v *ArticleView) Content() *Content {
	return v.content
}

// PostAt returns article's post time
func (v *ArticleView) PostAt() time.Time {
	return v.postAt
}

// LastEditedAt returns articles last edited time
func (v *ArticleView) LastEditedAt() time.Time {
	return v.lastEditedAt
}
