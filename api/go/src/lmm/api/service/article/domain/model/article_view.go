package model

import (
	"time"

	"lmm/api/model"
)

// ArticleView model defines what should be shown on article page
type ArticleView struct {
	model.ValueObject
	id           *ArticleID
	content      *Content
	postAt       time.Time
	lastEditedAt time.Time
}

// NewArticleView creates new ArticleView
func NewArticleView(articleID *ArticleID, content *Content, postAt, lastEditedAt time.Time) *ArticleView {
	return &ArticleView{
		id:           articleID,
		content:      content,
		postAt:       postAt,
		lastEditedAt: lastEditedAt,
	}
}

// ID returns article's id
func (v *ArticleView) ID() *ArticleID {
	return v.id
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
