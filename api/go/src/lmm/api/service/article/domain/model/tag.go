package model

import (
	"regexp"
	"strings"

	"lmm/api/service/article/domain"
)

var (
	patternTagName = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]{1,30}$")
)

// Tag is the tag model
type Tag struct {
	name  string
	order uint
}

// NewTag creates a new tag
func NewTag(name string, order uint) (*Tag, error) {
	name, err := validateTagName(name)
	if err != nil {
		return nil, err
	}

	return &Tag{name: name, order: order}, nil
}

// Name returns the tag's name
func (tag *Tag) Name() string {
	return tag.name
}

func (tag *Tag) Order() uint {
	return tag.order
}

func validateTagName(s string) (string, error) {
	name := strings.TrimSpace(s)
	if !patternTagName.MatchString(name) {
		return "", domain.ErrInvalidTagName
	}
	return name, nil
}

// TagView is a model used to view tag
type TagView struct {
	name  string
	count int
}

// NewTagView creates a new TagView
func NewTagView(name string, count int) *TagView {
	return &TagView{
		name:  name,
		count: count,
	}
}

// Name returns tag name
func (v *TagView) Name() string {
	return v.name
}

// Count returns the number of tag which named v.Name
func (v *TagView) Count() int {
	return v.count
}
