package model

import (
	"regexp"

	"lmm/api/context/article/domain"
	"lmm/api/domain/model"
)

var patternArticleID = regexp.MustCompile("[0-9a-zA-Z]{6}")

// ArticleID is the model to identify article
type ArticleID struct {
	model.ValueObject
	id string
}

// NewArticleID is a constructor of article id
func NewArticleID(s string) (*ArticleID, error) {
	id := ArticleID{}
	if err := id.setID(s); err != nil {
		return nil, err
	}
	return &id, nil
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
