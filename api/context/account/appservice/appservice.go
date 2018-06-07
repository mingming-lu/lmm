package appservice

import (
	"errors"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/db"
)

var (
	ErrDuplicateUserName         = errors.New("User name duplicated")
	ErrEmptyUserNameOrPassword   = errors.New("Empty user name or password")
	ErrInvalidInput              = errors.New("Invalid input")
	ErrInvalidToken              = errors.New("Invalid token")
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
	return &Usecase{repo: repo}
}

// SignIn is a usecase which users sign in with a account
func (uc *Usecase) SignIn(name, password string) (*model.User, error) {
	if name == "" || password == "" {
		return nil, ErrEmptyUserNameOrPassword
	}

	user, err := uc.repo.FindByName(name)
	if err != nil {
		if err.Error() == db.ErrNoRows.Error() {
			return nil, ErrInvalidUserNameOrPassword
		}
		return nil, err
	}

	if user.VerifyPassword(password) != nil {
		return nil, ErrInvalidUserNameOrPassword
	}

	return user, nil
}
