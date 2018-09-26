package appservice

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/service"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/context/account/infra"
	"lmm/api/testing"
)

func TestSignIn_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := randomUserNameAndPassword()

	_, err := app.SignUp(testing.StructToRequestBody(auth))
	t.NoError(err)

	user, err := app.SignIn(testing.StructToRequestBody(auth))
	t.NoError(err)
	t.Isa(&model.User{}, user)
	t.Is(auth.Name, user.Name())
	t.NoError(user.VerifyPassword(auth.Password))
}

func TestSignIn_InvalidPassword(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := randomUserNameAndPassword()

	_, err := app.SignUp(testing.StructToRequestBody(auth))
	t.NoError(err)

	auth.Password = "1234"
	user, err := app.SignIn(testing.StructToRequestBody(auth))

	t.Error(err)
	t.Nil(user)
	t.Is(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_NoSuchUser(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := randomUserNameAndPassword()
	user, err := app.SignIn(testing.StructToRequestBody(auth))

	t.Error(err)
	t.Nil(user)
	t.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_EmptyUserName(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := Auth{Name: "", Password: "1234"}

	res, err := app.SignIn(testing.StructToRequestBody(auth))

	t.Error(err)
	t.Nil(res)
	t.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_EmptyPassword(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := Auth{Name: "username", Password: ""}

	res, err := app.SignIn(testing.StructToRequestBody(auth))
	t.Error(err)
	t.Nil(res)
	t.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_EmptyUserNameAndPassword(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := Auth{Name: "", Password: ""}

	res, err := app.SignIn(testing.StructToRequestBody(auth))
	t.Error(err)
	t.Nil(res)
	t.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_Exception(tt *testing.T) {
	t := testing.NewTester(tt)
	app := New(testingService.NewMockedRepo())

	auth := randomUserNameAndPassword()

	_, err := app.SignUp(testing.StructToRequestBody(auth))
	t.Is("Cannot save user", err.Error())

	user, err := app.SignIn(testing.StructToRequestBody(auth))

	t.Error(err)
	t.Nil(user)
	t.Is("DB crashed", err.Error())
}
