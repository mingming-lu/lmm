package testing

import (
	"errors"
	"lmm/api/context/account/domain/model"
	"lmm/api/db"
	"lmm/api/domain/repository"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func NewUser() *model.User {
	db := db.New()
	defer db.Close()

	stmt1 := db.MustPrepare(`INSERT INTO user (name, password, guid, token, created_at) VALUES(?, ?, ?, ?, ?)`)
	defer stmt1.Close()

	name := uuid.New()[:32]
	password := uuid.New()
	user := model.NewUser(name, password)

	result, err := stmt1.Exec(user.Name, user.Password, user.GUID, user.Token, user.CreatedAt)
	if err != nil {
		panic(err)
	}
	userID, err := result.LastInsertId()

	stmt2 := db.MustPrepare(`SELECT id, name, guid, token, created_at FROM user WHERE id = ?`)
	defer stmt2.Close()

	user = &model.User{}
	if err := stmt2.QueryRow(userID).Scan(&user.ID, &user.Name, &user.GUID, &user.Token, &user.CreatedAt); err != nil {
		panic(err)
	}
	user.Password = password
	return user
}

type MockedRepo struct {
	repository.Repository
	testing.Mock
}

func NewMockedRepo() *MockedRepo {
	return &MockedRepo{Repository: &repository.Default{}}
}

func (repo *MockedRepo) FindByName(name string) (*model.User, error) {
	return nil, errors.New("DB crashed")
}

func (repo *MockedRepo) Put(*model.User) (*model.User, error) {
	return nil, errors.New("Cannot save user")
}

func (repo *MockedRepo) FindByToken(token string) (*model.User, error) {
	return nil, errors.New("No such user")
}
