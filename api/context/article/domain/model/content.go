package model

import "lmm/api/context/base/domain/model"

// Content shows article's content those are editable
type Content struct {
	model.ValueObject
	text *Text
	tags []*Tag
}

// NewContent returns a new Content value object pointer
func NewContent(text *Text, tags []*Tag) *Content {
	return &Content{text: text, tags: tags}
}

// Text returns article content's text
func (c *Content) Text() *Text {
	return c.text
}

// Tags returns article content's tags
func (c *Content) Tags() []*Tag {
	return c.tags
}
