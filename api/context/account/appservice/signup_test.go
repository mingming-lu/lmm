package appservice

import (
	"lmm/api/context/account/domain/repository"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
)

func TestSignUp(t *testing.T) {
	tester := testing.NewTester(t)
	repo := repository.New(testing.DB())
	app := New(repo)

	name, password := randomUserNameAndPassword()
	id, err := app.SignUp(name, password)
	tester.NoError(err)

	user, err := repo.FindByName(name)
	tester.NoError(err)
	tester.Is(id, user.ID())
	tester.Is(name, user.Name())
	tester.NoError(user.VerifyPassword(password))
}

func TestSignUp_Duplicate(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	name, password := randomUserNameAndPassword()
	id, err := app.SignUp(name, password)
	tester.NoError(err)
	tester.Not(uint64(0), id)

	id, err = app.SignUp(name, password)
	tester.Is(ErrDuplicateUserName, err)
	tester.Is(uint64(0), id)
}

func TestSignUp_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	id, err := app.SignUp("", "password")
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_EmptyPassword(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	id, err := app.SignUp("username", "")
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_EmptyUserNameAndPassword(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	id, err := app.SignUp("", "")
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_Exception(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(testingService.NewMockedRepo())

	name, password := randomUserNameAndPassword()
	id, err := app.SignUp(name, password)

	tester.Error(err)
	tester.Is(uint64(0), id)
	tester.Is("Cannot save user", err.Error())
}
