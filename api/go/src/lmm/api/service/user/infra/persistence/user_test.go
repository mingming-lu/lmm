package persistence

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"

	"lmm/api/clock"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/testing"
	transaction "lmm/api/transaction/datastore"
	"lmm/api/util/uuidutil"
)

func TestUserStore(tt *testing.T) {
	t := testing.NewTester(tt)

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if !t.NoError(err) {
		t.Fatalf(`failed to setup datastore: "%s"`, err.Error())
	}
	txManager := transaction.NewTransactionManager(client)

	repo := NewUserStore(client)

	tt.Run("Save", func(tt *testing.T) {
		user := newUser()

		tt.Run("Insert", func(tt *testing.T) {
			t := testing.NewTester(tt)

			t.NoError(txManager.RunInTransaction(ctx, func(c context.Context) error {
				return repo.Save(c, user)
			}))

			tt.Run("FindByName", func(tt *testing.T) {
				t := testing.NewTester(tt)

				testFindByName(ctx, t, txManager, repo, user.Name(), user)
			})
		})

		t.NoError(user.ChangeEmail(user.Name() + "@example.net"))

		tt.Run("Update", func(tt *testing.T) {
			t := testing.NewTester(tt)

			err := txManager.RunInTransaction(ctx, func(c context.Context) error {
				return repo.Save(c, user)
			})

			t.NoError(err)

			tt.Run("FindByName", func(tt *testing.T) {
				t := testing.NewTester(tt)

				testFindByName(ctx, t, txManager, repo, user.Name(), user)
			})
		})
	})

	tt.Run("FindByName", func(tt *testing.T) {
		user := newUser()

		tt.Run("NotFound", func(tt *testing.T) {
			t := testing.NewTester(tt)

			txManager.RunInReadOnly(ctx, func(c context.Context) error {
				found, err := repo.FindByName(c, user.Name())
				t.IsError(domain.ErrNoSuchUser, err)
				t.Nil(found)

				return nil
			})
		})
	})
}

func testFindByName(
	ctx context.Context,
	t *testing.Tester,
	txManager *transaction.TransactionManager,
	repo repository.UserRepository,
	name string,
	expect *model.User,
) bool {
	return txManager.RunInTransaction(ctx, func(c context.Context) error {
		found, err := repo.FindByName(c, name)
		if !t.NoError(err) {
			t.Fatal(err)
		}

		t.NotNil(found)

		t.Is(expect.Name(), found.Name())
		t.Is(expect.Email(), found.Email())
		t.Is(expect.Role(), found.Role())
		t.Is(expect.Password(), found.Password())
		t.Is(expect.Token(), found.Token())
		t.Is(expect.RegisteredAt(), found.RegisteredAt().UTC())

		return nil
	}) == nil
}

func newUser() *model.User {
	username := "u" + uuidutil.NewUUID()[:8]
	password := uuidutil.NewUUID()
	email := username + "@example.com"
	role := model.Admin
	token := uuidutil.NewUUID()
	registedAt := clock.Now()

	user, err := model.NewUser(username, email, password, token, role, registedAt)

	if err != nil {
		panic("internal error: " + err.Error())
	}

	return user
}
