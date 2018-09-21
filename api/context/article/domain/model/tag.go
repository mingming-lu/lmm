package model

import (
	"regexp"

	"lmm/api/context/article/domain"
	"lmm/api/domain/model"
	"lmm/api/utils/strings"
)

var (
	patternTagName = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]$")
)

// Tag is the tag model
type Tag struct {
	model.Entity
	id *TagID
}

// NewTag creates a new tag
func NewTag(articleID *ArticleID, name string) (*Tag, error) {
	name, err := validateTagName(name)
	if err != nil {
		return nil, err
	}

	id := &TagID{articleID: articleID, name: name}
	return &Tag{id: id}, nil
}

// ID returns the tag's id
func (tag *Tag) ID() *TagID {
	return tag.id
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
