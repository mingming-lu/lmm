package appservice

import (
	"lmm/api/context/account/domain/repository"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
)

func TestSignUp(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	id, err := New(repository.New()).SignUp(auth.Name, auth.Password)
	tester.NoError(err)
	tester.Is(uint64(1), id)
}

func TestSignUp_Duplicate(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	user := testingService.NewUser()

	repo := repository.New()
	id, err := New(repo).SignUp(user.Name, user.Password)
	tester.Error(err)
	tester.Is(ErrDuplicateUserName.Error(), err.Error())
	tester.Is(uint64(0), id)
}

func TestSignUp_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)

	id, err := New(repository.New()).SignUp("", "1234")
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_EmptyPassword(t *testing.T) {
	tester := testing.NewTester(t)

	id, err := New(repository.New()).SignUp("user", "")
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_EmptyUserNameAndPassword(t *testing.T) {
	tester := testing.NewTester(t)

	id, err := New(repository.New()).SignUp("", "")
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_Exception(t *testing.T) {
	tester := testing.NewTester(t)

	id, err := New(testingService.NewMockedRepo()).SignUp("foobar", "1234")

	tester.Error(err)
	tester.Is(uint64(0), id)
	tester.Is("Cannot save user", err.Error())
}
