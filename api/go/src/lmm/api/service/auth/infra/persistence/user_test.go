package persistence

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"lmm/api/service/auth/domain/service"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

func TestUserStorage(tt *testing.T) {
	c := context.Background()

	user := testutil.NewUser(dbEngine)

	tt.Run("FindByName", func(tt *testing.T) {
		tt.Run("Found", func(tt *testing.T) {
			t := testing.NewTester(tt)
			userFound, err := userRepo.FindByName(c, user.Name())
			t.NoError(err)
			t.NotNil(userFound)
			t.Is(user.Name(), userFound.Name())
			t.Is(user.RawToken(), userFound.RawToken())
		})

		tt.Run("NotFound", func(tt *testing.T) {
			t := testing.NewTester(tt)
			userFound, err := userRepo.FindByName(c, "whatever")
			t.IsError(sql.ErrNoRows, errors.Cause(err))
			t.Nil(userFound)
		})
	})

	tt.Run("FindByToken", func(tt *testing.T) {
		tt.Run("Found", func(tt *testing.T) {
			t := testing.NewTester(tt)
			token, err := service.NewTokenService().Encode(user.RawToken())
			if !t.NoError(err) {
				t.FailNow()
			}
			userFound, err := userRepo.FindByToken(c, token)
			t.NoError(err)
			t.NotNil(userFound)
			t.Is(user.Name(), userFound.Name())
			t.Is(token.Raw(), userFound.RawToken())
		})
		tt.Run("NotFound", func(tt *testing.T) {
			t := testing.NewTester(tt)
			otherToken, err := service.NewTokenService().Encode("whatever")
			if !t.NoError(err) {
				t.FailNow()
			}
			userFound, err := userRepo.FindByToken(c, otherToken)
			t.IsError(sql.ErrNoRows, errors.Cause(err))
			t.Nil(userFound)
		})
	})
}
