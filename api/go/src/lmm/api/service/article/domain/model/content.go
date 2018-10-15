package model

import (
	"sort"

	"lmm/api/domain/model"
	"lmm/api/service/article/domain"
)

// Content shows article's content those are editable
type Content struct {
	model.ValueObject
	text *Text
	tags []*Tag
}

// NewContent returns a new Content value object pointer
func NewContent(text *Text, tags []*Tag) (*Content, error) {
	c := Content{}

	if err := c.setText(text); err != nil {
		return nil, err
	}

	if err := c.setTags(tags); err != nil {
		return nil, err
	}

	return &c, nil
}

// Text returns article content's text
func (c *Content) Text() *Text {
	return c.text
}

func (c *Content) setText(text *Text) error {
	c.text = text
	return nil
}

// Tags returns article content's tags
func (c *Content) Tags() []*Tag {
	return c.tags
}

func (c *Content) setTags(tags []*Tag) error {
	var articleID *ArticleID
	if len(tags) > 0 {
		articleID = tags[0].ID().ArticleID()
	}

	for _, tag := range tags {
		if tag.ID().ArticleID().String() != articleID.String() {
			return domain.ErrTagsNotBelongToSameArticle
		}
	}

	lessTagOrder := func(i, j int) bool {
		return tags[i].ID().Order() < tags[j].ID().Order()
	}

	if !sort.SliceIsSorted(tags, lessTagOrder) {
		sort.Slice(tags, lessTagOrder)
	}

	c.tags = tags
	return nil
}
