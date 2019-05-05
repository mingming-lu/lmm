package datastore

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
)

const userKind = "user"

type userStore struct {
	datastore *datastore.Client
}

type user struct {
	Name      string    `datastore:"name"`
	Email     string    `datastore:"email"`
	Password  string    `datastore:"password,noindex"`
	Token     string    `datastore:"token"`
	Role      string    `datastore:"role,noindex"`
	CreatedAt time.Time `datastore:"created_at"`
}

// NewUserStore returns a UserRepository
func NewUserStore(client *datastore.Client) repository.UserRepository {
	return &userStore{datastore: client}
}

func (s *userStore) Save(c context.Context, userModel *model.User) error {
	id := datastore.NameKey(userKind, userModel.Name(), nil)

	_, err := s.datastore.RunInTransaction(c, func(tx *datastore.Transaction) error {
		usr := new(user)

		if err := tx.Get(id, usr); err != datastore.ErrNoSuchEntity {
			if err == nil {
				return domain.ErrUserNameAlreadyUsed
			}
			return err
		}

		_, err := tx.Put(id, &user{
			Name:      userModel.Name(),
			Email:     userModel.Email(),
			Password:  userModel.Password(),
			Token:     userModel.Token(),
			Role:      userModel.Role().Name(),
			CreatedAt: userModel.RegisteredAt(),
		})

		return err
	})

	return err
}

func (s *userStore) FindByName(c context.Context, username string) (*model.User, error) {
	id := datastore.NameKey(userKind, username, nil)

	usr := new(user)

	if err := s.datastore.Get(c, id, usr); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, domain.ErrNoSuchUser
		}
		return nil, err
	}

	return model.NewUser(
		usr.Name,
		usr.Email,
		usr.Password,
		usr.Token,
		model.NewRole(usr.Role),
		usr.CreatedAt,
	)
}

func (s *userStore) DescribeAll(c context.Context, options repository.DescribeAllOptions) ([]*model.UserDescriptor, uint, error) {
	panic("not implemented")
}
