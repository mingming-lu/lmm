package domain

import "errors"

var (
	ErrArticleTitleTooLong        = errors.New("article title too long")
	ErrEmptyArticleTitle          = errors.New("empty article title")
	ErrInvalidArticleID           = errors.New("invalid article id")
	ErrInvalidAliasArticleID      = errors.New("invalid alias article id")
	ErrInvalidArticleTitle        = errors.New("invalid article title")
	ErrInvalidTagName             = errors.New("invalid tag name")
	ErrNoSuchArticle              = errors.New("no such article")
	ErrNoSuchUser                 = errors.New("no such user")
	ErrNotArticleAuthor           = errors.New("only author allowed to edit article")
	ErrTagsNotBelongToSameArticle = errors.New("tags are not belong to same article")
)
