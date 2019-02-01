package persistence

import (
	"context"

	"github.com/pkg/errors"

	"lmm/api/clock"
	"lmm/api/http"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/service/user/domain/service"
	"lmm/api/storage/db"
	"lmm/api/util/mysqlutil"
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
		insert into user (name, password, token, role, created_at) values(?, ?, ?, ?, ?)
	`)
	defer stmt.Close()

	now := clock.Now()

	_, err := stmt.Exec(c, user.Name(), user.Password(), user.Token(), user.Role().Name(), now)

	if key, _, ok := mysqlutil.CheckDuplicateKeyError(err); ok && key == "name" {
		return errors.Wrap(domain.ErrUserNameAlreadyUsed, err.Error())
	}

	return err
}

// FindByName implementation
func (s *UserStorage) FindByName(c context.Context, username string) (*model.User, error) {
	panic("not implemented")
}

// DescribeByName implementation
func (s *UserStorage) DescribeByName(c context.Context, username string) (*model.UserDescriptor, error) {
	stmt := s.db.Prepare(c, `select role from user where name = ?`)
	defer stmt.Close()

	var rolename string

	if err := stmt.QueryRow(c, username).Scan(&rolename); err != nil {
		return nil, err
	}

	role := service.RoleAdapter(rolename)
	if role == model.Guest {
		http.Log().Panic(c, "expected not a guest")
	}

	return model.NewUserDescriptor(username, role)
}
