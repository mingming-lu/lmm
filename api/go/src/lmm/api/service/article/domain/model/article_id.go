package model

import (
	"regexp"
	"strings"

	"lmm/api/domain/model"
	"lmm/api/service/article/domain"
)

var patternArticleID = regexp.MustCompile("^[0-9a-z-]{32,64}$")

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

func (id *ArticleID) String() string {
	return id.id
}

func (id *ArticleID) setID(anID string) error {
	anID = strings.ToLower(anID)
	if !patternArticleID.MatchString(anID) {
		return domain.ErrInvalidArticleID
	}
	id.id = anID
	return nil
}
