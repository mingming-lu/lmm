package service

import (
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/storage"

	"github.com/akinaru-lu/errors"
)

var (
	ErrDuplicateUserName         = errors.New("user name duplicated")
	ErrEmptyUserNameOrPassword   = errors.New("empty user name or password")
	ErrInvalidAuthorization      = errors.New("invalid authorization")
	ErrInvalidInput              = errors.New("invalid input")
	ErrInvalidToken              = errors.New("invalid token")
	ErrInvalidUserNameOrPassword = errors.New("invalid user name or password")
	ErrInvalidBody               = errors.New("invalid body")
	ErrNoSuchUser                = errors.New("no such user")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name, password string) (*model.User, error) {
	user, err := factory.NewUser(name, password)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Add(user); err != nil {
		key, _, ok := storage.CheckErrorDuplicate(err.Error())
		if !ok {
			return nil, err
		}
		if key == "name" {
			return nil, ErrDuplicateUserName
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(name, password string) (*model.User, error) {
	user, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if err := user.VerifyPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByHashedToken(hashedToken string) (*model.User, error) {
	rawToken, err := DecodeToken(hashedToken)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByToken(rawToken)
}
