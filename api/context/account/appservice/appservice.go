package appservice

import (
	"errors"
	"lmm/api/context/account/domain/factory"
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

type AppService struct {
	repo repository.Repository
}

func New(repo repository.Repository) *AppService {
	return &Usecase{repo: repo}
}

func (app *AppService) SignUp(name, password string) (uint64, error) {
	if name == "" || password == "" {
		return 0, ErrEmptyUserNameOrPassword
	}

	m, err := factory.NewUser(name, password)
	if err != nil {
		return 0, err
	}
	err = uc.repo.Add(m)
	if err != nil {
		key, _, ok := repository.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "name" {
			return 0, ErrDuplicateUserName
		}
		return 0, err
	}
	return user.ID(), nil
}

// SignIn is a usecase which users sign in with a account
func (app *AppService) SignIn(name, password string) (*model.User, error) {
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
