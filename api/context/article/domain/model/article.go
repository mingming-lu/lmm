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
	id     *ArticleID
	text   *Text
	author *Author
	tags   []*Tag
}

// NewArticle is a article constructor
func NewArticle(articleID *ArticleID, text *Text, author *Author, tags []*Tag) *Article {
	return &Article{
		id:     articleID,
		text:   text,
		author: author,
		tags:   tags,
	}
}

// ID returns the id of the article
func (a *Article) ID() *ArticleID {
	return a.id
}

// Text returns the text of the article
func (a *Article) Text() *Text {
	return a.text
}

// Author returns the author of the article
func (a *Article) Author() *Author {
	return a.author
}

// Tags returns all tags of the article
func (a *Article) Tags() []*Tag {
	return a.tags
}

// EditText changes the text of the article
func (a *Article) EditText(newText *Text) error {
	if a.Text().Equals(newText) {
		return ErrArticleTextNoChange
	}
	a.text = newText
	return nil
}

// AddTag adds a new tag to article
func (a *Article) AddTag(tag *Tag) error {
	for _, t := range a.tags {
		if tag.Equals(t) {
			return ErrTagAlreadyAdded
		}
	}
	a.tags = append(a.tags, tag)
	return nil
}

// RemoveTag removes a tag from the article
func (a *Article) RemoveTag(tag *Tag) error {
	targetIndex := -1
	for index, t := range a.Tags() {
		if tag.Equals(t) {
			targetIndex = index
			break
		}
	}
	if targetIndex == -1 {
		return ErrTagNotExists
	}
	a.tags = append(a.tags[:targetIndex], a.tags[targetIndex+1:]...)
	return nil
}
