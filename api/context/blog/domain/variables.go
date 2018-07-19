package domain

import "errors"

// blog
var (
	ErrBlogNoChange        = errors.New("blog no change")
	ErrBlogTitleDuplicated = errors.New("blog title duplicated")
	ErrCategoryNotSet      = errors.New("category not set")
	ErrEmptyBlogTitle      = errors.New("blog title cannot be empty")
	ErrInvalidBlogID       = errors.New("invalid blog id")
	ErrInvalidCount        = errors.New("invalid count")
	ErrInvalidPage         = errors.New("invalid page")
	ErrNoPermission        = errors.New("no permission")
	ErrNoSuchBlog          = errors.New("no such blog")
)

// category
var (
	ErrCategoryNoChanged     = errors.New("category no changed")
	ErrDuplicateCategoryName = errors.New("duplicate category name")
	ErrInvalidCategoryID     = errors.New("invalid category id")
	ErrInvalidCategoryName   = errors.New("invalid category name")
	ErrNoSuchCategory        = errors.New("no such category")
)

// tag
var (
	ErrInvalidTagName   = errors.New("invalid tag name")
	ErrDuplicateTagName = errors.New("duplicate tag name")
)
