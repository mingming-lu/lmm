package persistence

import (
	"context"
	"time"

	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/storage/db"
)

// UserStorage implements user/domain/repository.UserRepository
type UserStorage struct {
	db db.DB
}

// NewUserStorage returns a UserRepository
func NewUserStorage(db db.DB) repository.UserRepository {
	return &UserStorage{db: db}
}

// Save persists a user model
func (s *UserStorage) Save(c context.Context, user *model.User) error {
	stmt := s.db.Prepare(c, `
		insert into user (name, password, token, created_at) values(?, ?, ?, ?)
	`)
	defer stmt.Close()

	now := time.Now()

	_, err := stmt.Exec(c, user.Name(), user.Password(), user.Token(), now)
	// TODO check user name conflict

	return err
}

// FindByUserName implementation
func (s *UserStorage) FindByUserName(c context.Context, username string) (*model.User, error) {
	panic("not implemented")
}
