package ui

import (
	"io"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/factory"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"strings"
)

func TestPostV1SignIn_400_InvalidInput(t *testing.T) {
	res := postSignIn(strings.NewReader("not a json"))

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(http.Status(http.StatusBadRequest)+"\n", res.Body())
}

func TestPostV1SignIn_400_EmptyUserName(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: "1234"})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_400_EmptyPassword(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "foobar", Password: ""})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_400_EmptyUserNameAndPassword(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: ""})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidUserName(t *testing.T) {
	name, password := uuid.New()[:32], uuid.New()
	factory.NewUser(name, password)

	requestBody := testing.StructToRequestBody(Auth{Name: "name", Password: password})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(appservice.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidPassword(t *testing.T) {
	name, password := uuid.New()[:32], uuid.New()
	factory.NewUser(name, password)

	requestBody := testing.StructToRequestBody(Auth{Name: name, Password: "1234"})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(appservice.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidUserNameAndPassword(t *testing.T) {
	name, password := uuid.New()[:32], uuid.New()
	requestBody := testing.StructToRequestBody(Auth{Name: name, Password: password})
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
