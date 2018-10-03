package domain

import "github.com/pkg/errors"

var (
	// ErrNoSuchUser error
	ErrNoSuchUser = errors.New("no such user")
)
