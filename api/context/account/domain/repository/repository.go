package repository

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/db"
	"lmm/api/domain/repository"

	"github.com/akinaru-lu/errors"
)

type Repository interface {
	repository.Repository
	Put(*model.User) (*model.User, error)
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
func (repo *repo) Put(user *model.User) (*model.User, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`INSERT INTO user (name, password, guid, token, created_at) VALUES (?, ?, ?, ?, ?)`)
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Password, user.GUID, user.Token, user.CreatedAt.UTC())
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = uint64(id)
	return user, nil
}

// FindByName return a user model determined by name
func (repo *repo) FindByName(name string) (*model.User, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name, password, guid, token, created_at FROM user WHERE name = ?`)
	defer stmt.Close()

	user := &model.User{}
	err := stmt.QueryRow(name).Scan(&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return user, nil
}

func (repo *repo) FindByToken(token string) (*model.User, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name, password, guid, token, created_at FROM user WHERE token = ?`)
	defer stmt.Close()

	user := model.User{}
	err := stmt.QueryRow(token).Scan(&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &user, nil
}

func FindByToken(token string) (*model.User, error) {
	db := db.New()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name, password, guid, token, created_at FROM user WHERE token = ?`)
	defer stmt.Close()

	user := model.User{}
	err := stmt.QueryRow(token).Scan(&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &user, nil
}
