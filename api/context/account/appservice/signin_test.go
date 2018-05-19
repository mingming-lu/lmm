package appservice

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
)

func TestSignIn_Success(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)
	uc := New(repository.New())

	m := testingService.NewUser()
	user, err := uc.SignIn(m.Name, m.Password)

	tester.NoError(err)
	tester.Is(m.ID, user.ID)
	tester.Is(m.Name, user.Name)
	tester.Isa(&model.User{}, user)
}

func TestSignIn_InvalidPassword(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)
	uc := New(repository.New())

	m := testingService.NewUser()
	user, err := uc.SignIn(m.Name, "1234")

	tester.Error(err)
	tester.Nil(user)
	tester.Is(ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_NoSuchUser(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)
	uc := New(repository.New())

	res, err := uc.SignIn("foobar", "1234")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrInvalidUserNameOrPassword, err)
}

func TestSignIn_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)
	uc := New(repository.New())

	res, err := uc.SignIn("", "1234")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignIn_EmptyPassword(t *testing.T) {
	tester := testing.NewTester(t)
	uc := New(repository.New())

	res, err := uc.SignIn("user", "")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignIn_EmptyUserNameAndPassword(t *testing.T) {
	tester := testing.NewTester(t)
	uc := New(repository.New())

	res, err := uc.SignIn("", "")
	tester.Error(err)
	tester.Nil(res)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignIn_Exception(t *testing.T) {
	tester := testing.NewTester(t)

	user, err := New(testingService.NewMockedRepo()).SignIn("foobar", "1234")

	tester.Error(err)
	tester.Nil(user)
	tester.Is("DB crashed", err.Error())
}
