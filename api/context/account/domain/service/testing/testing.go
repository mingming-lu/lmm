package testing

import (
	"errors"
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/model"
	"lmm/api/db"
	"lmm/api/domain/repository"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"time"
)

func NewUser() *model.User {
	db := db.New()
	defer db.Close()

	stmt1 := db.MustPrepare(`INSERT INTO user (id, name, password, token, created_at) VALUES(?, ?, ?, ?, ?)`)
	defer stmt1.Close()

	name := uuid.New()[:32]
	password := uuid.New()
	user, err := factory.NewUser(name, password)
	panic(err)

	result, err := stmt1.Exec(user.ID(), user.Name(), user.Password(), user.Token(), user.CreatedAt())
	if err != nil {
		panic(err)
	}
	userID, err := result.LastInsertId()

	stmt2 := db.MustPrepare(`SELECT name, password, token, created_at FROM user WHERE id = ?`)
	defer stmt2.Close()

	var (
		userName      string
		userToken     string
		userCreatedAt time.Time
	)

	if err := stmt2.QueryRow(userID).Scan(&userName, &userToken, &userCreatedAt); err != nil {
		panic(err)
	}
	return user
}

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
