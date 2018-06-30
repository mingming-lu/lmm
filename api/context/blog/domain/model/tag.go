package model

import (
	"errors"
	"lmm/api/domain/model"
	"lmm/api/utils/strings"
	"regexp"
)

var (
	ErrInvalidTagName = errors.New("invalid tag name")
)

var (
	patternValidTagName = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]{1,31}$")
)

type Tag struct {
	model.Entity
	id   uint64
	data tagData
}

type tagData struct {
	model.ValueObject
	blogID  uint64
	tagName tagName
}

type tagName string

func newTagName(s string) (tagName, error) {
	name := strings.TrimSpace(s)
	if !patternValidCategoryName.MatchString(name) {
		return "", ErrInvalidTagName
	}
	return tagName(name), nil
}

func NewTag(id, blogID uint64, name string) (*Tag, error) {
	s, err := newTagName(name)
	if err != nil {
		return nil, err
	}
	data := tagData{blogID: blogID, tagName: s}
	return &Tag{id: id, data: data}, nil
}

func (t *Tag) ID() uint64 {
	return t.id
}

func (t *Tag) BlogID() uint64 {
	return t.data.blogID
}

func (t *Tag) Name() string {
	return string(t.data.tagName)
}

func (t *Tag) UpdateName(newName string) error {
	name, err := newTagName(newName)
	if err != nil {
		return err
	}
	t.data.tagName = name
	return nil
}
