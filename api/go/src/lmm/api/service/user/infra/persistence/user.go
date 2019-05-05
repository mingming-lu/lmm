package persistence

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	transaction "lmm/api/transaction/datastore"
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
	tx, err := transaction.FromContext(c)
	if err != nil {
		return err
	}

	id := datastore.NameKey(userKind, userModel.Name(), nil)

	if _, err := tx.Mutate(datastore.NewUpsert(id, &user{
		Name:      userModel.Name(),
		Email:     userModel.Email(),
		Password:  userModel.Password(),
		Token:     userModel.Token(),
		Role:      userModel.Role().Name(),
		CreatedAt: userModel.RegisteredAt(),
	})); err != nil {
		return err
	}

	return nil
}

func (s *userStore) FindByName(c context.Context, username string) (*model.User, error) {
	tx, err := transaction.FromContext(c)
	if err != nil {
		return nil, err
	}

	usr := new(user)

	if err := tx.Get(datastore.NameKey(userKind, username, nil), usr); err != nil {
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
