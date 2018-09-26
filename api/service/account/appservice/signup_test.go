package appservice

import (
	"lmm/api/context/account/domain/service"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/context/account/infra"
	"lmm/api/testing"
)

func TestSignUp(t *testing.T) {
	tester := testing.NewTester(t)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := randomUserNameAndPassword()
	id, err := app.SignUp(testing.StructToRequestBody(auth))
	tester.NoError(err)

	user, err := repo.FindByName(auth.Name)
	tester.NoError(err)
	tester.Is(id, user.ID())
	tester.Is(auth.Name, user.Name())
	tester.NoError(user.VerifyPassword(auth.Password))
}

func TestSignUp_Duplicate(t *testing.T) {
	tester := testing.NewTester(t)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := randomUserNameAndPassword()
	id, err := app.SignUp(testing.StructToRequestBody(auth))
	tester.NoError(err)
	tester.Not(uint64(0), id)

	id, err = app.SignUp(testing.StructToRequestBody(auth))
	tester.IsError(service.ErrDuplicateUserName, err)
	tester.Is(uint64(0), id)
}

func TestSignUp_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := Auth{Name: "", Password: "password"}
	id, err := app.SignUp(testing.StructToRequestBody(auth))
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignUp_EmptyPassword(t *testing.T) {
	tester := testing.NewTester(t)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := Auth{Name: "username", Password: ""}

	id, err := app.SignUp(testing.StructToRequestBody(auth))
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignUp_EmptyUserNameAndPassword(t *testing.T) {
	tester := testing.NewTester(t)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := Auth{Name: "", Password: ""}

	id, err := app.SignUp(testing.StructToRequestBody(auth))
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.IsError(service.ErrInvalidUserNameOrPassword, err)
}

func TestSignUp_Exception(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(testingService.NewMockedRepo())

	auth := randomUserNameAndPassword()

	id, err := app.SignUp(testing.StructToRequestBody(auth))

	tester.Error(err)
	tester.Is(uint64(0), id)
	tester.Is("Cannot save user", err.Error())
}
