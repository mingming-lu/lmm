package model

import (
	"regexp"
)

var (
	patternAliasArticleID = regexp.MustCompile("^[0-9a-z-]{8,80}$")
)

// ArticleID is the model to identify article
type ArticleID int64

// NewArticleID is a constructor of article id
func NewArticleID(raw int64) *ArticleID {
	id := ArticleID(raw)
	return &id
}
