package model

import (
	"errors"
	"time"

	"lmm/api/context/base/domain/model"
)

var (
	ErrArticleTextNoChange = errors.New("article text no change")
	ErrTagAlreadyAdded     = errors.New("the tag has already been added")
	ErrTagNotAdded         = errors.New("tag not added")
)

type Article struct {
	model.Entity
	id             ArticleID
	text           ArticleText
	writer         ArticleWriter
	postAt         time.Time
	lastModifiedAt time.Time
	tags           []*Tag
}

func NewArticle(
	articleID ArticleID,
	text ArticleText,
	writer ArticleWriter,
	postAt time.Time,
	lastModifiedAt time.Time,
	tags []*Tag,
) *Article {
	return &Article{
		id:             articleID,
		text:           text,
		writer:         writer,
		postAt:         postAt,
		lastModifiedAt: lastModifiedAt,
		tags:           tags,
	}
}

func (a *Article) ID() ArticleID {
	return a.id
}

func (a *Article) Text() ArticleText {
	return a.text
}

func (a *Article) Writer() ArticleWriter {
	return a.writer
}

func (a *Article) PostAt() time.Time {
	return a.postAt
}

func (a *Article) LastModifiedAt() time.Time {
	return a.lastModifiedAt
}

func (a *Article) Tags() []*Tag {
	return a.tags
}

func (a *Article) setLastModifiedAt(at time.Time) {
	a.lastModifiedAt = at
}

func (a *Article) ModifyText(newText ArticleText) error {
	if a.Text().Equals(newText) {
		return ErrArticleTextNoChange
	}
	a.setLastModifiedAt(time.Now())
	a.text = newText
	return nil
}

func (a *Article) AddTag(tag *Tag) error {
	for _, t := range a.tags {
		if tag.Equals(t) {
			return ErrTagAlreadyAdded
		}
	}
	a.tags = append(a.tags, tag)
	return nil
}

func (a *Article) RemoveTag(tag *Tag) error {
	targetIndex := -1
	for index, t := range a.Tags() {
		if tag.Equals(t) {
			targetIndex = index
			break
		}
	}
	if targetIndex == -1 {
		return ErrTagNotAdded
	}
	a.tags = append(a.tags[:targetIndex], a.tags[targetIndex+1:]...)
	return nil
}
