package domain

import "errors"

var (
	ErrDuplicateImageID = errors.New("duplicate image id")
	ErrEmptyImageType   = errors.New("no specified image type")
	ErrFailedToUpload   = errors.New("failed to upload image")
	ErrInvalidCount     = errors.New("invalid count")
	ErrInvalidPage      = errors.New("invalid page")
	ErrMarkImageFailed  = errors.New("failed to mark image's type")
	ErrNoSuchImage      = errors.New("no such image")
	ErrNoSuchImageType  = errors.New("no such image type")
)
