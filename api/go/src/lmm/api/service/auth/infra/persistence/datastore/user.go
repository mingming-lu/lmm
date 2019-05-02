package datastore

import (
	"context"

	"cloud.google.com/go/datastore"

	"golang.org/x/xerrors"

	"lmm/api/service/auth/domain/model"
	"lmm/api/service/auth/domain/repository"
)

const userKind = "user"

type userStore struct {
	datastore *datastore.Client
}

type user struct {
	Name     string `datastore:"name"`
	Password string `datastore:"password"`
	Token    string `datastore:"token"`
	Role     string `datastore:"role"`
}

// NewUserStore creates new UserStorage
func NewUserStore(client *datastore.Client) repository.UserRepository {
	return &userStore{datastore: client}
}

func (s *userStore) FindByName(ctx context.Context, name string) (*model.User, error) {
	user, err := s.findByQuery(ctx, datastore.NewQuery(userKind).Filter("name =", name))
	if err != nil {
		return nil, xerrors.Errorf(`cannot find user named "%s": %w`, name, err)
	}
	return user, nil
}

func (s *userStore) FindByToken(ctx context.Context, token *model.Token) (*model.User, error) {
	user, err := s.findByQuery(ctx, datastore.NewQuery(userKind).Filter("token =", token.Raw()))
	if err != nil {
		return nil, xerrors.Errorf(`cannot find user by token: %w`, err)
	}
	return user, nil
}

func (s *userStore) findByQuery(ctx context.Context, q *datastore.Query) (*model.User, error) {
	usr := new(user)

	iter := s.datastore.Run(ctx, q)

	if _, err := iter.Next(usr); err != nil {
		return nil, err
	}

	return model.NewUser(usr.Name, usr.Password, usr.Token, usr.Role), nil
}
