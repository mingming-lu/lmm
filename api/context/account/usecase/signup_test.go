package usecase

import (
	"lmm/api/context/account/domain/repository"
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

	auth := Auth{Name: "foobar", Password: "1234"}
	repo := repository.New()
	New(repo).SignUp(auth.Name, auth.Password)
	id, err := New(repo).SignUp(auth.Name, auth.Password)
	tester.Error(err)
	tester.Is(ErrDuplicateUserName, err)
	tester.Is(uint64(0), id)
}
