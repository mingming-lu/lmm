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

	stmt := db.Must(`INSERT INTO user (name, password, guid, token, created_at) VALUES (?, ?, ?, ?, ?)`)
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

	stmt := db.Must(`SELECT id, name, password, guid, token, created_at FROM user WHERE name = ?`)
	defer stmt.Close()

	user := &model.User{}
	err := stmt.QueryRow(name).Scan(&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return user, nil
}

func Add(name, password, guid, token string) (int64, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must(`INSERT INTO user (name, password, guid, token) values (?, ?, ?, ?)`)
	defer stmt.Close()

	res, err := stmt.Exec(name, password, guid, token)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func ByToken(token string) (*model.User, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT id, name, password, guid, token, created_at FROM user WHERE token = ?")
	defer stmt.Close()

	user := model.User{}
	err := stmt.QueryRow(token).Scan(
		&user.ID, &user.Name, &user.Password, &user.GUID, &user.Token, &user.CreatedAt,
	)

	return &user, err
}
