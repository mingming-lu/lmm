package domain

import "errors"

var (
	ErrArticleTitleTooLong = errors.New("article title too long")
	ErrEmptyArticleTitle   = errors.New("empty article title")
	ErrInvalidArticleID    = errors.New("invalid article id")
	ErrInvalidArticleTitle = errors.New("invalid article title")
	ErrInvalidCount        = errors.New("invalid count")
	ErrInvalidPage         = errors.New("invalid page")
	ErrInvalidTagName      = errors.New("invalid tag name")
	ErrNoSuchArticle       = errors.New("no such article")
)
