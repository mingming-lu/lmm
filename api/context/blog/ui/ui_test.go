package ui

import (
	accountFactory "lmm/api/context/account/domain/factory"
	accountModel "lmm/api/context/account/domain/model"
	accountService "lmm/api/context/account/domain/service"
	accountStorage "lmm/api/context/account/infra"
	account "lmm/api/context/account/ui"
	"lmm/api/context/blog/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *accountModel.User
var ui *UI
var accountUI *account.UI

func TestMain(m *testing.M) {
	initUser()
	initUI()

	accountUI = account.New(accountStorage.NewUserStorage(testing.DB()))

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

	err = accountStorage.NewUserStorage(testing.DB()).Add(user)
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

func initUI() {
	db := testing.DB()
	blogRepo := infra.NewBlogStorage(db)
	categoryRepo := infra.NewCategoryStorage(db)
	ui = New(blogRepo, categoryRepo)
}
