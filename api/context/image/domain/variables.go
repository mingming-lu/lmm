package domain

import "errors"

var (
	ErrDuplicateImageID = errors.New("duplicate image id")
	ErrNoSuchImage      = errors.New("no such image")
)
