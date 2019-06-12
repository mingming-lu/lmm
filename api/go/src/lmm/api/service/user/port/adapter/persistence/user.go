package persistence

import (
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

type user struct {
	ID           *datastore.Key `datastore:"__key__"`
	Name         string         `datastore:"Name"`
	Email        string         `datastore:"Email"`
	Password     string         `datastore:"Password,noindex"`
	Token        string         `datastore:"Token"`
	Role         string         `datastore:"Role,noindex"`
	RegisteredAt time.Time      `datastore:"RegisteredAt,noindex"`
}

const (
	userKind = "User"
)

// UserDataStore implements UserRepository
type UserDataStore struct {
	source *datastore.Client
	transaction.Manager
}

func NewUserDataStore(source *datastore.Client) *UserDataStore {
	return &UserDataStore{
		source:  source,
		Manager: dsUtil.NewTransactionManager(source),
	}
}

func (s *UserDataStore) NextID(tx transaction.Transaction) (model.UserID, error) {
	keys, err := s.source.AllocateIDs(tx, []*datastore.Key{
		datastore.IncompleteKey(userKind, nil),
	})

	if err != nil {
		return -1, errors.Wrap(err, "failed to allocate new id")
	}

	return model.UserID(keys[0].ID), nil
}

// Save implementation
func (s *UserDataStore) Save(tx transaction.Transaction, model *model.User) error {
	k := datastore.IDKey(userKind, int64(model.ID()), nil)

	_, err := dsUtil.MustTransaction(tx).Mutate(
		datastore.NewUpsert(k, &user{
			ID:           k,
			Name:         model.Name(),
			Email:        model.Email(),
			Password:     model.Password(),
			Token:        model.Token(),
			Role:         model.Role().Name(),
			RegisteredAt: model.RegisteredAt(),
		}),
	)

	return errors.Wrap(err, "faile to save user to datastore")
}

func (s *UserDataStore) findByFilter(tx transaction.Transaction, filter, value string) (*model.User, error) {
	q := datastore.NewQuery(userKind).KeysOnly().Filter(filter, value).Limit(1)

	keys, err := s.source.GetAll(tx, q, nil)
	if err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchUser, err.Error())
	}

	if len(keys) == 0 {
		return nil, domain.ErrNoSuchUser
	}

	var user user
	if err := dsUtil.MustTransaction(tx).Get(keys[0], &user); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, domain.ErrNoSuchUser
		}
		return nil, errors.Wrap(err, "internal error: failed to get user by key")
	}

	return model.NewUser(
		model.UserID(user.ID.ID),
		user.Name,
		user.Email,
		user.Password,
		user.Token,
		model.RoleFromString(user.Role),
		user.RegisteredAt,
	)
}

// FindByName implementation
func (s *UserDataStore) FindByName(tx transaction.Transaction, username string) (*model.User, error) {
	return s.findByFilter(tx, "Name =", username)
}

// FindByToken implementation
func (s *UserDataStore) FindByToken(tx transaction.Transaction, token string) (*model.User, error) {
	return s.findByFilter(tx, "Token =", token)
}
