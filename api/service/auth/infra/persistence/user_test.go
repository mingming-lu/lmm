package persistence

import (
	"context"
	"lmm/api/util/testingutil"

	"github.com/pkg/errors"

	"lmm/api/service/auth/domain"
	"lmm/api/service/auth/domain/service"
	"lmm/api/testing"
)

func TestUserStorage(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	user, err := testingutil.NewAuthUser(dbEngine)
	if !t.NoError(err) {
		t.FailNow()
	}

	token, err := service.NewTokenService().Encode(user.RawToken())
	if !t.NoError(err) {
		t.FailNow()
	}

	t.Run("FindByName", func(_ *testing.T) {
		t.Run("Found", func(_ *testing.T) {
			userFound, err := userRepo.FindByName(c, user.Name())
			t.NoError(err)
			t.NotNil(userFound)
			t.Is(user.Name(), userFound.Name())
			t.Is(token.Raw(), userFound.RawToken())
		})

		t.Run("NotFound", func(_ *testing.T) {
			userFound, err := userRepo.FindByName(c, "whatever")
			t.IsError(domain.ErrNoSuchUser, errors.Cause(err))
			t.Nil(userFound)
		})
	})

	t.Run("FindByToken", func(_ *testing.T) {
		t.Run("Found", func(_ *testing.T) {
			userFound, err := userRepo.FindByToken(c, token)
			t.NoError(err)
			t.NotNil(userFound)
			t.Is(user.Name(), userFound.Name())
			t.Is(token.Raw(), userFound.RawToken())
		})
		t.Run("NotFound", func(_ *testing.T) {
			otherToken, err := service.NewTokenService().Encode("whatever")
			if !t.NoError(err) {
				t.FailNow()
			}
			userFound, err := userRepo.FindByToken(c, otherToken)
			t.IsError(domain.ErrNoSuchUser, errors.Cause(err))
			t.Nil(userFound)
		})
	})
}
