package repository

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/domain/repository"
	"time"

	"github.com/akinaru-lu/errors"
)

type Repository interface {
	repository.Repository
	Add(*model.User) error
	FindByName(string) (*model.User, error)
	FindByToken(string) (*model.User, error)
}

type repo struct {
	repository.Default
}

func New() Repository {
	return new(repo)
}

// Put puts a new user into repository and return a User model with generated id
func (repo *repo) Add(user *model.User) error {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`INSERT INTO user (id, name, password, token, created_at) VALUES (?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(user.Name(), user.Password(), user.Token(), user.CreatedAt().UTC())
	if err != nil {
		return err
	}
	return nil
}

// FindByName return a user model determined by name
func (repo *repo) FindByName(name string) (*model.User, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name, password, token, created_at FROM user WHERE name = ?`)
	defer stmt.Close()

	var (
		userID        uint64
		userName      string
		userPassword  string
		userToken     string
		userCreatedAt time.Time
	)
	err := stmt.QueryRow(name).Scan(&userID, &userName, &userPassword, &userToken, &userCreatedAt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return model.NewUser(userID, userName, userPassword, userToken, userCreatedAt), nil
}

func (repo *repo) FindByToken(token string) (*model.User, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name, password, guid, token, created_at FROM user WHERE token = ?`)
	defer stmt.Close()

	var (
		userID        uint64
		userName      string
		userPassword  string
		userToken     string
		userCreatedAt time.Time
	)
	err := stmt.QueryRow(token).Scan(&userID, &userName, &userPassword, &userToken, &userCreatedAt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return model.NewUser(userID, userName, userPassword, userToken, userCreatedAt), nil
}
