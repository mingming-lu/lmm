package domain

import "github.com/pkg/errors"

var (
	// ErrNoPermission error
	ErrNoPermission = errors.New("no permission")

	// ErrInvalidPassword error
	ErrInvalidPassword = errors.New("invalid password")

	// ErrInvalidUserID error
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrInvalidUserName error
	ErrInvalidUserName = errors.New("invalid user name")

	// ErrInvalidUserToken error
	ErrInvalidUserToken = errors.New("invalid user token, expect a uuid v4")

	// ErrInvalidPermission error
	ErrInvalidPermission = errors.New("invalid permission")

	// ErrNoSuchUser error
	ErrNoSuchUser = errors.New("no such user")

	// ErrNoSuchRole error
	ErrNoSuchRole = errors.New("no such role")

	// ErrUserNameAlreadyUsed error
	ErrUserNameAlreadyUsed = errors.New("user name has already been used")

	// ErrUserPasswordEmpty error
	ErrUserPasswordEmpty = errors.New("user password is empty")

	// ErrUserPasswordTooLong error
	ErrUserPasswordTooLong = errors.New("user password should be equal to or shorter than 250")

	// ErrUserPasswordTooShort error
	ErrUserPasswordTooShort = errors.New("user password should be equal to or longer than 8")

	// ErrUserPasswordTooWeak error
	ErrUserPasswordTooWeak = errors.New("user password is too weak")
)
