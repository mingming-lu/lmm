package usecase

import (
	"lmm/api/context/account/domain/repository"
	"lmm/api/testing"
)

func TestSignIn_Success(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)
	uc := New(repository.New())

	id, err := uc.SignUp("foobar", "1234")
	tester.NoError(err)

	res, err := uc.SignIn("foobar", "1234")
	tester.NoError(err)
	tester.Is(id, res.ID)
	tester.Is("foobar", res.Name)
}

func TestSignIn_InvalidPassword(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)
	uc := New(repository.New())

	_, err := uc.SignUp("foobar", "1234")
	tester.NoError(err)

	res, err := uc.SignIn("foobar", "5555")
	tester.Error(err)
	tester.Nil(res)
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
