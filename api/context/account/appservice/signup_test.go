package appservice

import (
	"lmm/api/context/account/domain/repository"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestSignUp(t *testing.T) {
	tester := testing.NewTester(t)
	repo := repository.New()

	name, password := uuid.New()[:32], uuid.New()
	id, err := New(repo).SignUp(name, password)
	tester.NoError(err)

	user, err := repo.FindByName(name)
	tester.NoError(err)
	tester.Is(id, user.ID())
	tester.Is(name, user.Name())
	tester.NoError(user.VerifyPassword(password))
}

func TestSignUp_Duplicate(t *testing.T) {
	tester := testing.NewTester(t)
	repo := repository.New()

	name, password := uuid.New()[:32], uuid.New()
	id, err := New(repo).SignUp(name, password)
	tester.NoError(err)
	tester.Not(uint64(0), id)

	id, err = New(repo).SignUp(name, password)
	tester.Is(ErrDuplicateUserName, err)
	tester.Is(uint64(0), id)
}

func TestSignUp_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)
	repo := repository.New()

	id, err := New(repo).SignUp("", uuid.New())
	tester.Error(err)
	tester.Is(uint64(id), id)
	tester.Is(ErrEmptyUserNameOrPassword, err)
}

func TestSignUp_EmptyPassword(t *testing.T) {
	tester := testing.NewTester(t)

	id, err := New(repository.New()).SignUp(uuid.New()[:32], "")
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
