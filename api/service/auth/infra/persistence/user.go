package persistence

import (
	"context"

	"lmm/api/service/auth/domain"
	"lmm/api/service/auth/domain/model"
	"lmm/api/service/auth/domain/repository"
	"lmm/api/storage/db"
)

// UserStorage is a UserRepository implementation
type UserStorage struct {
	db db.DB
}

// NewUserStorage creates new UserStorage
func NewUserStorage(db db.DB) repository.UserRepository {
	return &UserStorage{db: db}
}

// FindByName implementation
func (s *UserStorage) FindByName(c context.Context, name string) (*model.User, error) {
	return s.findUser(c, `select name, password, token from user where name = ?`, name)
}

// FindByToken implementation
func (s *UserStorage) FindByToken(c context.Context, token *model.Token) (*model.User, error) {
	return s.findUser(c, `select name, password, token from user where token = ?`, token.Raw())
}

func (s *UserStorage) findUser(c context.Context, query string, args ...interface{}) (*model.User, error) {
	stmt := s.db.Prepare(c, query)
	defer stmt.Close()

	var (
		userName     string
		userPassword string
		userToken    string
	)

	err := stmt.QueryRow(c, args...).Scan(&userName, &userPassword, &userToken)
	if err != nil {
		if err == db.ErrNoRows {
			return nil, domain.ErrNoSuchUser
		}
		return nil, err
	}

	return model.NewUser(userName, userPassword, userToken), nil
}
