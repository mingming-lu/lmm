package domain

import "github.com/pkg/errors"

var (
	ErrInvalidAccessToken      = errors.New("invalid access token auth format")
	ErrInvalidBasicAuthFormat  = errors.New("invalid basic auth format")
	ErrInvalidBearerAuthFormat = errors.New("invalid bearer auth format")
	ErrInvalidGrantType        = errors.New("invalid grant type")
	ErrInvalidAuthToken        = errors.New("invalid auth token")
	ErrInvalidTokenFormat      = errors.New("invalid token format")
	ErrInvalidTokenLength      = errors.New("invalid token length")
	ErrNoSuchUser              = errors.New("no such user")
	ErrPasswordNotMatched      = errors.New("wrong password")
	ErrTokenExpired            = errors.New("token expired")
)

const (
	GrantTypeBasicAuth    = "basicAuth"
	GrantTypeRefreshToken = "refreshToken"
)
