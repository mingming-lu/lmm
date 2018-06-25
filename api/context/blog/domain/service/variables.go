package service

import "errors"

// blog
var (
	ErrInvalidBlogID  = errors.New("invalid blog id")
	ErrCategoryNotSet = errors.New("category not set")
)

// category
var (
	ErrInvalidCategoryID = errors.New("invalid category id")
	ErrNoSuchCategory    = errors.New("no such category")
)
