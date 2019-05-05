package datastore

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/datastore"

	"lmm/api/clock"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/testing"
	"lmm/api/util/uuidutil"
)

func TestUserStore(tt *testing.T) {
	t := testing.NewTester(tt)

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if !t.NoError(err) {
		t.Fatalf(`failed to setup datastore: "%s"`, err.Error())
	}

	repo := NewUserStore(client)

	tt.Run("Save", func(tt *testing.T) {
		tt.Run("Serialize", func(tt *testing.T) {
			user := newUser()

			tt.Run("Success", func(tt *testing.T) {
				t := testing.NewTester(tt)

				t.NoError(repo.Save(ctx, user))
			})

			tt.Run("Duplicate", func(tt *testing.T) {
				t := testing.NewTester(tt)

				t.IsError(domain.ErrUserNameAlreadyUsed, repo.Save(ctx, user))
			})

			tt.Run("FindByName", func(tt *testing.T) {
				t := testing.NewTester(tt)

				found, err := repo.FindByName(ctx, user.Name())
				t.NoError(err)
				t.Is(user.Name(), found.Name())
				t.Is(user.Email(), found.Email())
				t.Is(user.Role(), found.Role())
				t.Is(user.Password(), found.Password())
				t.Is(user.Token(), found.Token())
				t.Is(user.RegisteredAt(), found.RegisteredAt().UTC())
			})
		})

		tt.Run("Concurrently", func(tt *testing.T) {
			t := testing.NewTester(tt)

			user := newUser()

			errorCounter := struct {
				count int
				sync.Mutex
				sync.WaitGroup
			}{}

			conum := 16

			for i := 0; i < conum; i++ {
				errorCounter.Add(1)
				go func() {
					if repo.Save(ctx, user) != nil {
						errorCounter.Lock()
						errorCounter.count++
						errorCounter.Unlock()
					}
					errorCounter.Done()
				}()
			}

			errorCounter.Wait()

			t.Is(conum-1, errorCounter.count)

			tt.Run("FindByName", func(tt *testing.T) {
				t := testing.NewTester(tt)

				found, err := repo.FindByName(ctx, user.Name())
				t.NoError(err)
				t.NoError(err)
				t.Is(user.Name(), found.Name())
				t.Is(user.Email(), found.Email())
				t.Is(user.Role(), found.Role())
				t.Is(user.Password(), found.Password())
				t.Is(user.Token(), found.Token())
				t.Is(user.RegisteredAt(), found.RegisteredAt().UTC())
			})
		})
	})

	tt.Run("FindByName", func(tt *testing.T) {
		user := newUser()

		tt.Run("NotFound", func(tt *testing.T) {
			t := testing.NewTester(tt)

			found, err := repo.FindByName(ctx, user.Name())
			t.IsError(domain.ErrNoSuchUser, err)
			t.Nil(found)
		})
	})
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
