package model

import (
	"lmm/api/context/article/domain"
	"lmm/api/utils/strings"

	"regexp"
)

var (
	patternTagName = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]$")
)

type TagID struct {
	articleID ArticleID
	name      string
}

func (id TagID) ArticleID() ArticleID {
	return id.articleID
}

func (id TagID) Name() string {
	return id.name
}

type Tag struct {
	id TagID
}

func NewTag(articleID ArticleID, name string) (*Tag, error) {
	name, err := validateTagName(name)
	if err != nil {
		return nil, err
	}

	id := TagID{articleID: articleID, name: name}
	return &Tag{id: id}, nil
}

func (tag *Tag) ID() TagID {
	return tag.id
}

func validateTagName(s string) (string, error) {
	name := strings.TrimSpace(s)
	if !patternTagName.MatchString(name) {
		return "", domain.ErrInvalidTagName
	}
	return name, nil
}
