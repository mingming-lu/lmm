package datastore

import (
	"context"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/datastore"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/service"
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
		})
	})
}

func newUser() *model.User {
	username := "u" + uuidutil.NewUUID()[:8]
	password := uuidutil.NewUUID()
	email := username + "@example.com"
	role := service.RoleAdapter("admin")
	token := uuidutil.NewUUID()
	registedAt := time.Now()

	user, err := model.NewUser(username, email, password, token, role, registedAt)

	if err != nil {
		panic("internal error: " + err.Error())
	}

	return user
}
