package ui

import (
	"io"
	"lmm/api/context/account/appservice"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
	"net/http"
	"strings"
)

func TestPostV1SignIn_400_InvalidInput(t *testing.T) {
	testing.InitTable("user")

	res := postSignIn(strings.NewReader("not a json"))

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(http.StatusText(http.StatusBadRequest)+"\n", res.Body())
}

func TestPostV1SignIn_400_EmptyUserName(t *testing.T) {
	testing.InitTable("user")

	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: "1234"})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_400_EmptyPassword(t *testing.T) {
	testing.InitTable("user")

	requestBody := testing.StructToRequestBody(Auth{Name: "foobar", Password: ""})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_400_EmptyUserNameAndPassword(t *testing.T) {
	testing.InitTable("user")

	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: ""})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidUserName(t *testing.T) {
	testing.InitTable("user")
	user := testingService.NewUser()

	requestBody := testing.StructToRequestBody(Auth{Name: user.Name, Password: "123"})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(appservice.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidPassword(t *testing.T) {
	testing.InitTable("user")
	user := testingService.NewUser()

	requestBody := testing.StructToRequestBody(Auth{Name: "123", Password: user.Password})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(appservice.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidUserNameAndPassword(t *testing.T) {
	testing.InitTable("user")
	testingService.NewUser()

	requestBody := testing.StructToRequestBody(Auth{Name: "123", Password: "123"})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(appservice.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func postSignIn(requestBody io.Reader) *testing.Response {
	res := testing.NewResponse()

	router := testing.NewRouter()
	router.POST("/v1/signin", SignIn)
	router.ServeHTTP(res, testing.POST("/v1/signin", requestBody))

	return res
}
