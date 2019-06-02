package model

import (
	"regexp"
	"strings"

	"lmm/api/service/article/domain"
	"lmm/api/service/base/model"
)

var (
	patternArticleID      = regexp.MustCompile("^[0-9a-z-]{8,64}$")
	patternAliasArticleID = regexp.MustCompile("^[0-9a-z-]{8,80}$")
)

// ArticleID is the model to identify article
type ArticleID struct {
	model.ValueObject
	id    string
	alias string
}

// NewArticleID is a constructor of article id
func NewArticleID(s string) (*ArticleID, error) {
	id := ArticleID{}
	if err := id.setID(s); err != nil {
		return nil, err
	}
	return &id, nil
}

// Raw gets raw article id in string
func (id *ArticleID) Raw() string {
	return id.id
}

// String gets alias or raw article id
func (id *ArticleID) String() string {
	if id.alias != "" {
		return id.alias
	}
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

// SetAlias sets alias article id
func (id *ArticleID) SetAlias(alias string) error {
	alias = strings.ToLower(alias)
	if !patternAliasArticleID.MatchString(alias) {
		return domain.ErrInvalidAliasArticleID
	}
	if alias == id.id {
		alias = ""
		return nil
	}
	id.alias = alias
	return nil
}