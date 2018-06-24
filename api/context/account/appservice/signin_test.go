package appservice

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
)

func TestSignIn_Success(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	name, password := randomUserNameAndPassword()
	app.SignUp(name, password)

	user, err := app.SignIn(name, password)
	tester.NoError(err)
	tester.Isa(&model.User{}, user)
	tester.Is(name, user.Name())
	tester.NoError(user.VerifyPassword(password))
}

func TestSignIn_InvalidPassword(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	name, password := randomUserNameAndPassword()
	app.SignUp(name, password)

	user, err := app.SignIn(name, "1234")
	tester.Error(err)
	tester.Nil(user)
	tester.Is(ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_NoSuchUser(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	name, password := randomUserNameAndPassword()
	user, err := app.SignIn(name, password)

	tester.Error(err)
	tester.Nil(user)
	tester.Is(ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)
	uc := New(repository.New(testing.DB()))

	res, err := uc.SignIn("", "1234")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignIn_EmptyPassword(t *testing.T) {
	tester := testing.NewTester(t)
	uc := New(repository.New(testing.DB()))

	res, err := uc.SignIn("username", "")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignIn_EmptyUserNameAndPassword(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	res, err := app.SignIn("", "")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignIn_Exception(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(testingService.NewMockedRepo())

	name, password := randomUserNameAndPassword()
	app.SignUp(name, password)
	user, err := app.SignIn(name, password)

	tester.Error(err)
	tester.Nil(user)
	tester.Is("DB crashed", err.Error())
}
