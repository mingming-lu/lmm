package testing

import (
	"errors"
	"lmm/api/context/account/domain/model"
	"lmm/api/domain/repository"
	"lmm/api/testing"
)

type MockedRepo struct {
	repository.Repository
	testing.Mock
}

func NewMockedRepo() *MockedRepo {
	return &MockedRepo{Repository: &repository.Default{}}
}

func (repo *MockedRepo) Add(*model.User) error {
	return errors.New("Cannot save user")
}

func (repo *MockedRepo) FindByName(name string) (*model.User, error) {
	return nil, errors.New("DB crashed")
}

func (repo *MockedRepo) FindByToken(token string) (*model.User, error) {
	return nil, errors.New("No such user")
}
