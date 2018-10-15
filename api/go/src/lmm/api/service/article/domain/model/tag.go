package model

import (
	"regexp"
	"strings"

	"lmm/api/domain/model"
	"lmm/api/service/article/domain"
)

var (
	patternTagName = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]{1,30}$")
)

// Tag is the tag model
type Tag struct {
	model.Entity
	id   *TagID
	name string
}

// NewTag creates a new tag
func NewTag(articleID *ArticleID, order uint, name string) (*Tag, error) {
	name, err := validateTagName(name)
	if err != nil {
		return nil, err
	}

	id := &TagID{articleID: articleID, order: order}
	return &Tag{id: id, name: name}, nil
}

// ID returns the tag's id
func (tag *Tag) ID() *TagID {
	return tag.id
}

// Name returns the tag's name
func (tag *Tag) Name() string {
	return tag.name
}

// Equals compares tag with another
func (tag *Tag) Equals(anotherTag *Tag) bool {
	return tag.ID().Equals(anotherTag.ID())
}

func validateTagName(s string) (string, error) {
	name := strings.TrimSpace(s)
	if !patternTagName.MatchString(name) {
		return "", domain.ErrInvalidTagName
	}
	return name, nil
}
