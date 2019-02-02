package persistence

import (
	"context"
	"time"

	"github.com/pkg/errors"

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

	_, err := stmt.Exec(c, user.Name(), user.Password(), user.Token(), user.Role().Name(), user.RegisteredAt())

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
	stmt := s.db.Prepare(c, `select role, created_at from user where name = ?`)
	defer stmt.Close()

	var (
		rolename  string
		createdAt time.Time
	)

	if err := stmt.QueryRow(c, username).Scan(&rolename, &createdAt); err != nil {
		return nil, err
	}

	role := service.RoleAdapter(rolename)
	if role == model.Guest {
		http.Log().Panic(c, "expected not a guest")
	}

	return model.NewUserDescriptor(username, role, createdAt)
}

func (s *UserStorage) DescribeAll(c context.Context, options repository.DescribeAllOptions) ([]*model.UserDescriptor, error) {
	stmt := s.db.Prepare(c,
		`select name, role, created_at from user order by `+s.mappingOrder(options.Order)+` limit ? offset ?`)
	defer stmt.Close()

	rows, err := stmt.Query(c, options.Count, (options.Page-1)*options.Count)
	if err != nil {
		return nil, err
	}

	users := make([]*model.UserDescriptor, 0)

	var (
		username  string
		rolename  string
		createdAt time.Time
	)

	for rows.Next() {
		if err := rows.Scan(&username, &rolename, &createdAt); err != nil {
			return nil, err
		}
		role := service.RoleAdapter(rolename)
		user, err := model.NewUserDescriptor(username, role, createdAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	rows.Close()

	return users, nil
}

func (s *UserStorage) mappingOrder(order repository.DescribeAllOrder) string {
	switch order {
	case repository.DescribeAllOrderByNameAsc:
		return "name asc"
	case repository.DescribeAllOrderByNameDesc:
		return "name desc"
	case repository.DescribeAllOrderByRoleAsc:
		return "role asc, name asc"
	case repository.DescribeAllOrderByRoleDesc:
		return "role desc, name asc"
	case repository.DescribeAllOrderByRegisteredDateAsc:
		return "created_at asc, name asc"
	case repository.DescribeAllOrderByRegisteredDateDesc:
		return "created_at desc, name asc"
	default:
		panic("invalid order")
	}
}
