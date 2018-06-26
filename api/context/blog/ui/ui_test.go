package ui

import (
	accountFactory "lmm/api/context/account/domain/factory"
	accountModel "lmm/api/context/account/domain/model"
	accountRepository "lmm/api/context/account/domain/repository"
	accountService "lmm/api/context/account/domain/service"
	account "lmm/api/context/account/ui"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *accountModel.User
var ui *UI
var accountUI *account.UI

func TestMain(m *testing.M) {
	initUser()
	ui = New(testing.DB())
	accountUI = account.New(testing.DB())

	code := m.Run()
	os.Exit(code)
}

func initUser() {
	var err error

	name, password := uuid.New()[:31], uuid.New()
	user, err = accountFactory.NewUser(name, password)
	if err != nil {
		panic(err)
	}

	err = accountRepository.New(testing.DB()).Add(user)
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
}
