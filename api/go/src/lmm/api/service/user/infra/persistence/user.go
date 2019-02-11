package persistence

import (
	"context"
	"database/sql"
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
		insert into user (name, email, password, token, role, created_at) values(?, ?, ?, ?, ?, ?)
	`)
	defer stmt.Close()

	_, err := stmt.Exec(c, user.Name(), user.Email(), user.Password(), user.Token(), user.Role().Name(), user.RegisteredAt())

	if key, _, ok := mysqlutil.CheckDuplicateKeyError(err); ok && key == "name" {
		return errors.Wrap(domain.ErrUserNameAlreadyUsed, err.Error())
	}

	return err
}

// FindByName implementation
func (s *UserStorage) FindByName(c context.Context, username string) (*model.User, error) {
	stmt := s.db.Prepare(c, `select name, email, password, token, role, created_at from user where name = ?`)
	defer stmt.Close()

	var (
		name         string
		email        string
		password     string
		token        string
		rawRole      string
		registeredAt time.Time
	)

	if err := stmt.QueryRow(c, username).Scan(&name, &email, &password, &token, &rawRole, &registeredAt); err != nil {
		return nil, err
	}

	role := service.RoleAdapter(rawRole)

	return model.NewUser(name, email, password, token, role, registeredAt)
}

// DescribeAll implementation
func (s *UserStorage) DescribeAll(c context.Context, options repository.DescribeAllOptions) ([]*model.UserDescriptor, uint, error) {
	tx, err := s.db.Begin(c, &sql.TxOptions{
		ReadOnly:  true,
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return nil, 0, err
	}

	countUsers := tx.Prepare(c, `select count(*) from user`)
	defer countUsers.Close()

	selectUsers := tx.Prepare(c,
		`select name, email, role, created_at from user order by `+s.mappingOrder(options.Order)+` limit ? offset ?`)
	defer selectUsers.Close()

	var totalUsers uint
	if err := countUsers.QueryRow(c).Scan(&totalUsers); err != nil {
		return nil, 0, db.RollbackWithError(tx, err)
	}

	rows, err := selectUsers.Query(c, options.Count, (options.Page-1)*options.Count)
	if err != nil {
		return nil, 0, db.RollbackWithError(tx, err)
	}

	users := make([]*model.UserDescriptor, 0)

	var (
		username  string
		emailaddr string
		rolename  string
		createdAt time.Time
	)

	for rows.Next() {
		if err := rows.Scan(&username, &emailaddr, &rolename, &createdAt); err != nil {
			return nil, 0, db.RollbackWithError(tx, err)
		}
		role := service.RoleAdapter(rolename)
		user, err := model.NewUserDescriptor(username, emailaddr, role, createdAt)
		if err != nil {
			return nil, 0, db.RollbackWithError(tx, err)
		}
		users = append(users, user)
	}
	if err := rows.Close(); err != nil {
		return nil, 0, db.RollbackWithError(tx, err)
	}

	if err := tx.Commit(); err != nil {
		http.Log().Warn(c, err.Error())
	}

	return users, totalUsers, err
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
