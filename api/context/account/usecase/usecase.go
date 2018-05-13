package usecase

import (
	"errors"
	"lmm/api/context/account/domain/repository"
)

var (
	ErrInvalidInput      = errors.New("Invalid input")
	ErrDuplicateUserName = errors.New("User name duplicated")
	ErrIncorrectPassword = errors.New("Incorrect password")
)

type Auth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Usecase struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}
