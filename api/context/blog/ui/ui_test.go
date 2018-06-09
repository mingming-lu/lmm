package ui

import (
	accountFactory "lmm/api/context/account/domain/factory"
	accountModel "lmm/api/context/account/domain/model"
	accountRepository "lmm/api/context/account/domain/repository"
	accountService "lmm/api/context/account/domain/service"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *accountModel.User

func TestMain(m *testing.M) {
	var err error

	name, password := uuid.New()[:31], uuid.New()
	user, err = accountFactory.NewUser(name, password)
	if err != nil {
		panic(err)
	}

	err = accountRepository.New().Add(user)
	if err != nil {
		panic(err)
	}

	user = accountModel.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		accountService.EncodeToken(user.Token()),
		user.CreatedAt(),
	)

	code := m.Run()
	os.Exit(code)
}
