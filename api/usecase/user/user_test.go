package user

import (
	repo "lmm/api/domain/repository/user"
	"lmm/api/testing"
)

func TestSignUp(t *testing.T) {
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	requestBody := testing.StructToRequestBody(auth)

	testing.InitTable("user")
	id, err := New(repo.New()).SignUp(requestBody)
	tester.NoError(err)
	tester.Is(uint64(1), id)
}
