package datastore

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"

	"github.com/google/uuid"

	"lmm/api/service/auth/domain/service"
	"lmm/api/testing"
)

func TestUserStore(tt *testing.T) {
	t := testing.NewTester(tt)

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT_ID"))
	if !t.NoError(err) {
		t.Fatalf(`failed to setup datastore: "%s"`, err.Error())
	}

	username := "username"
	password := "password"
	token := uuid.New().String()
	role := "admin"

	key, err := client.Put(ctx, datastore.IncompleteKey(userKind, nil), &user{
		Name:     username,
		Password: password,
		Token:    token,
		Role:     role,
	})

	if !t.NoError(err) {
		t.Fatalf("failed to put user entity: %s", err.Error())
	}
	t.Logf("key: %#v", key)

	repo := NewUserStore(client)

	tt.Run("FindByName", func(tt *testing.T) {
		tt.Run("Found", func(tt *testing.T) {
			t := testing.NewTester(tt)

			user, err := repo.FindByName(ctx, username)
			t.NoError(err)
			t.Is(username, user.Name())

			t.Logf("%#v", user)
		})

		tt.Run("NotFound", func(tt *testing.T) {
			t := testing.NewTester(tt)

			user, err := repo.FindByName(ctx, "nosuchuser")
			t.Error(err)
			t.Nil(user)
			t.Log(err.Error())
		})
	})

	tt.Run("FindByToken", func(tt *testing.T) {
		tt.Run("Found", func(tt *testing.T) {
			t := testing.NewTester(tt)

			tokenModel, err := service.NewTokenService().Encode(token)

			user, err := repo.FindByToken(ctx, tokenModel)
			t.NoError(err)
			t.Is(token, user.RawToken())

			t.Logf("%#v", user)
		})

		tt.Run("NotFound", func(tt *testing.T) {
			t := testing.NewTester(tt)

			tokenModel, err := service.NewTokenService().Encode(uuid.New().String())

			user, err := repo.FindByToken(ctx, tokenModel)
			t.Error(err)
			t.Nil(user)
			t.Log(err.Error())
		})
	})
}
