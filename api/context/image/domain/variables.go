package domain

import "errors"

var (
	ErrDuplicateImageID = errors.New("duplicate image id")
	ErrNoSuchImage      = errors.New("no such image")
	ErrFailedToUpload   = errors.New("failed to upload image")
	ErrInvalidCount     = errors.New("invalid count")
	ErrInvalidPage      = errors.New("invalid page")
)
