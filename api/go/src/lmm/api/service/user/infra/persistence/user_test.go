package persistence

import (
	"context"
	"testing"
	"time"

	"lmm/api/clock"
	"lmm/api/pkg/transaction"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/infra/service"
	"lmm/api/util/uuidutil"

	"cloud.google.com/go/datastore"
	"github.com/stretchr/testify/assert"
)

func mustRandomUser(userDataStore *UserDataStore) *model.User {
	username := "U" + uuidutil.NewUUID()[:8]
	email := username + "@lmm.local"
	password := uuidutil.NewUUID()

	f := factory.NewFactory(&service.BcryptService{}, userDataStore)

	var user *model.User

	err := userDataStore.RunInTransaction(context.Background(), func(tx transaction.Transaction) error {
		userCreated, err := f.NewUser(tx, username, email, password)
		user = userCreated
		return err
	}, nil)

	if err != nil {
		panic("failed to create new user: " + err.Error())
	}

	return user
}

type testClock struct{}

func (c testClock) Now() time.Time {
	return time.Now().Truncate(1 * time.Second)
}

func TestUserDataStore(t *testing.T) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clock.DefaultClock = &testClock{}

	dataStore, err := datastore.NewClient(c, "")
	if err != nil {
		t.Fatal("failed to connect to datastore")
	}
	defer dataStore.Close()

	userDataStore := NewUserDataStore(dataStore)

	user := mustRandomUser(userDataStore)

	t.Run("NextID", func(t *testing.T) {
		userDataStore.RunInTransaction(c, func(tx transaction.Transaction) error {
			userID, err := userDataStore.NextID(tx)
			assert.NoError(t, err)
			assert.NotZero(t, userID)
			return nil
		}, nil)
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("Insert", func(t *testing.T) {
			userDataStore.RunInTransaction(c, func(tx transaction.Transaction) error {
				assert.NoError(t, userDataStore.Save(tx, user))
				return nil
			}, nil)
		})
	})

	t.Run("FindByName", func(t *testing.T) {
		t.Run("Found", func(t *testing.T) {
			userDataStore.RunInTransaction(c, func(tx transaction.Transaction) error {
				userFound, err := userDataStore.FindByName(tx, user.Name())
				if !assert.NoError(t, err) {
					t.Fatal(err)
				}

				assert.EqualValues(t, user, userFound)
				return nil
			}, nil)

			t.Run("ReadOnly", func(t *testing.T) {
				userDataStore.RunInTransaction(c, func(tx transaction.Transaction) error {
					userFound, err := userDataStore.FindByName(tx, user.Name())
					if !assert.NoError(t, err) {
						t.Fatal(err)
					}

					assert.EqualValues(t, user, userFound)
					return nil
				}, &transaction.Option{ReadOnly: true})
			})
		})

		t.Run("NotFound", func(t *testing.T) {
			user := mustRandomUser(userDataStore)

			userDataStore.RunInTransaction(c, func(tx transaction.Transaction) error {
				userFound, err := userDataStore.FindByName(tx, user.Name())
				assert.Error(t, domain.ErrNoSuchUser, err)
				assert.Nil(t, userFound)

				return nil
			}, nil)

			t.Run("ReadOnly", func(t *testing.T) {
				userDataStore.RunInTransaction(c, func(tx transaction.Transaction) error {
					userFound, err := userDataStore.FindByName(tx, user.Name())
					assert.Error(t, domain.ErrNoSuchUser, err)
					assert.Nil(t, userFound)

					return nil
				}, &transaction.Option{ReadOnly: true})
			})
		})
	})
}
