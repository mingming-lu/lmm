package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	repo "lmm/api/domain/repository/user"
	"lmm/api/testing"
)

func TestSignup(t *testing.T) {
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	b, err := json.Marshal(auth)
	tester.NoError(err)
	requestBody := ioutil.NopCloser(bytes.NewReader(b))

	testing.InitTable("user")
	id, err := New(repo.New()).SignUp(requestBody)
	tester.NoError(err)
	tester.Is(id, uint64(1))
}
