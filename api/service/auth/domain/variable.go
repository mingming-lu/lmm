package domain

import "github.com/pkg/errors"

var (
	ErrInvalidBasicAuthFormat  = errors.New("invalid basic auth format")
	ErrInvalidBearerAuthFormat = errors.New("invalid bearer auth format")
	ErrInvalidTokenFormat      = errors.New("invalid token format")
	ErrInvalidTokenLength      = errors.New("invalid token length")
	ErrNoSuchUser              = errors.New("no such user")
	ErrPasswordNotMatched      = errors.New("wrong password")
)
