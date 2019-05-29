package persistence

import (
	"context"
	"time"

	"lmm/api/pkg/transaction"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/service"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
)

type user struct {
	ID           *datastore.Key `datastore:"__key__"`
	Name         string         `datastore:"Name"`
	Email        string         `datastore:"Email"`
	Password     string         `datastore:"Password"`
	Token        string         `datastore:"Token"`
	Role         string         `datastore:"Role"`
	RegisteredAt time.Time      `datastore:"RegisteredAt"`
}

const (
	userKind = "User"
)

// UserDataStore implements UserRepository
type UserDataStore struct {
	source *datastore.Client
}

func NewUserDataStore(source *datastore.Client) *UserDataStore {
	return &UserDataStore{source: source}
}

type txImpl struct {
	context.Context
	*datastore.Transaction
}

func (tx *txImpl) Commit() error {
	_, err := tx.Transaction.Commit()
	return err
}

func (s *UserDataStore) Begin(c context.Context, opts *transaction.Option) (transaction.Transaction, error) {
	if opts != nil && opts.ReadOnly {
		tx, err := s.source.NewTransaction(c, datastore.ReadOnly)
		if err != nil {
			return nil, err
		}
		return &txImpl{c, tx}, nil
	}

	tx, err := s.source.NewTransaction(c)
	if err != nil {
		return nil, err
	}
	return &txImpl{c, tx}, nil
}

func (s *UserDataStore) RunInTransaction(c context.Context, f func(tx transaction.Transaction) error, opts *transaction.Option) error {
	tx, err := s.Begin(c, opts)
	if err != nil {
		return err
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			tx.Rollback()
			panic(recovered)
		}
	}()

	if err := f(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return errors.Wrap(err2, err.Error())
		}
		return err
	}

	return tx.Commit()
}

func mustTx(t transaction.Transaction) *txImpl {
	tx, ok := t.(*txImpl)
	if !ok {
		panic("not a *datastore.Transaction")
	}
	return tx
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

	_, err := mustTx(tx).Mutate(
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

// FindByName implementation
func (s *UserDataStore) FindByName(tx transaction.Transaction, username string) (*model.User, error) {
	q := datastore.NewQuery(userKind).KeysOnly().Filter("Name =", username).Limit(1)

	keys, err := s.source.GetAll(tx, q, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get key from user name: %s", username)
	}

	if len(keys) == 0 {
		return nil, domain.ErrNoSuchUser
	}

	var user user
	if err := mustTx(tx).Get(keys[0], &user); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, domain.ErrNoSuchUser
		}
		return nil, errors.Wrapf(err, "unexpected error when find user named %s", username)
	}

	return model.NewUser(
		model.UserID(user.ID.ID),
		user.Name,
		user.Email,
		user.Password,
		user.Token,
		service.RoleAdapter(user.Role),
		user.RegisteredAt,
	)
}

func (repo *UserDataStore) FindByToken(tx transaction.Transaction, token string) (*model.User, error) {
	panic("not implemented")
}
