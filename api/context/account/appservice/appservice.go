package appservice

import (
	"errors"
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/domain/service"
	"lmm/api/db"
	repoUtil "lmm/api/domain/repository"
)

var (
	ErrDuplicateUserName         = errors.New("User name duplicated")
	ErrEmptyUserNameOrPassword   = errors.New("Empty user name or password")
	ErrInvalidInput              = errors.New("Invalid input")
	ErrInvalidToken              = errors.New("Invalid token")
	ErrInvalidUserNameOrPassword = errors.New("Invalid user name or password")
)

type AppService struct {
	repo repository.Repository
}

func New(repo repository.Repository) *AppService {
	return &AppService{repo: repo}
}

func (app *AppService) SignUp(name, password string) (uint64, error) {
	if name == "" || password == "" {
		return 0, ErrEmptyUserNameOrPassword
	}

	user, err := factory.NewUser(name, password)
	if err != nil {
		return 0, err
	}
	err = app.repo.Add(user)
	if err != nil {
		key, _, ok := repoUtil.CheckErrorDuplicate(err.Error())
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

	user, err := app.repo.FindByName(name)
	if err != nil {
		if err.Error() == db.ErrNoRows.Error() {
			return nil, ErrInvalidUserNameOrPassword
		}
		return nil, err
	}

	if user.VerifyPassword(password) != nil {
		return nil, ErrInvalidUserNameOrPassword
	}

	return model.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		service.EncodeToken(user.Token()),
		user.CreatedAt(),
	), nil
}

func (app *AppService) VerifyToken(encodedToken string) (user *model.User, err error) {
	token, err := service.DecodeToken(encodedToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err = app.repo.FindByToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return user, nil
}
