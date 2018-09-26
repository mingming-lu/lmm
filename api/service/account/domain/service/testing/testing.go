package testing

import (
	"errors"
	"lmm/api/context/account/domain/model"
	"lmm/api/testing"
)

type MockedRepo struct {
	testing.Mock
}

func NewMockedRepo() *MockedRepo {
	return &MockedRepo{}
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
