package model

import (
	"regexp"

	"lmm/api/context/article/domain"
	"lmm/api/domain/model"
)

var patternArticleID = regexp.MustCompile("[0-9a-Z]{6}")

type ArticleID struct {
	model.ValueObject
	id string
}

func NewArticleID(s string) (ArticleID, error) {
	id := ArticleID{}
	if err := id.setID(s); err != nil {
		return id, err
	}
	return id, nil
}

func (id ArticleID) String() string {
	return id.id
}

func (id ArticleID) setID(anID string) error {
	if !patternArticleID.MatchString(anID) {
		return domain.ErrInvalidArticleID
	}
	id.id = anID
	return nil
}
