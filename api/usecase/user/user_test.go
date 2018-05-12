package user

import (
	repo "lmm/api/domain/repository/user"
	"lmm/api/testing"
)

func TestSignUp(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	requestBody := testing.StructToRequestBody(auth)

	id, err := New(repo.New()).SignUp(requestBody)
	tester.NoError(err)
	tester.Is(uint64(1), id)
}

func TestSignUp_Duplicate(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	New(repo.New()).SignUp(testing.StructToRequestBody(auth))
	id, err := New(repo.New()).SignUp(testing.StructToRequestBody(auth))
	tester.Error(err)
	tester.Is(ErrDuplicateUserName, err)
	tester.Is(uint64(0), id)
}
