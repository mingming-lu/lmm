package domain

import "errors"

var (
	ErrArticleTitleTooLong = errors.New("article title too long")
	ErrEmptyArticleTitle   = errors.New("empty article title")
	ErrInvalidArticleID    = errors.New("invalid article id")
	ErrInvalidArticleTitle = errors.New("invalid article title")
	ErrInvalidTagName      = errors.New("invalid tag name")
)
