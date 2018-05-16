package usecase

import (
	"errors"
	"lmm/api/context/account/domain/repository"
)

var (
	ErrDuplicateUserName         = errors.New("User name duplicated")
	ErrEmptyUserNameOrPassword   = errors.New("Empty user name or password")
	ErrInvalidInput              = errors.New("Invalid input")
	ErrInvalidUserNameOrPassword = errors.New("Invalid user name or password")
)

type Auth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Usecase struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Usecase {
	return &Usecase{repo: repository.New()}
}
