package ui

import (
	"io"
	"strings"

	"lmm/api/http"
	"lmm/api/service/account/domain/factory"
	"lmm/api/service/account/domain/service"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestPostV1SignIn_400_InvalidInput(tt *testing.T) {
	res := postSignIn(strings.NewReader("not a json"))

	t := testing.NewTester(tt)
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(http.Status(http.StatusBadRequest)+"\n", res.Body())
}

func TestPostV1SignIn_404_EmptyUserName(tt *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: "1234"})
	res := postSignIn(requestBody)

	t := testing.NewTester(tt)
	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_EmptyPassword(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "foobar", Password: ""})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_EmptyUserNameAndPassword(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: ""})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidUserName(t *testing.T) {
	name, password := uuid.New()[:31], uuid.New()
	factory.NewUser(name, password)

	requestBody := testing.StructToRequestBody(Auth{Name: "name", Password: password})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidPassword(t *testing.T) {
	name, password := uuid.New()[:31], uuid.New()
	factory.NewUser(name, password)

	requestBody := testing.StructToRequestBody(Auth{Name: name, Password: "1234"})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignIn_404_InvalidUserNameAndPassword(t *testing.T) {
	name, password := uuid.New()[:31], uuid.New()
	requestBody := testing.StructToRequestBody(Auth{Name: name, Password: password})
	res := postSignIn(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusNotFound, res.StatusCode())
	tester.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func postSignIn(requestBody io.Reader) *testing.Response {
	res := testing.NewResponse()

	router := testing.NewRouter()
	router.POST("/v1/signin", ui.SignIn)
	router.ServeHTTP(res, testing.POST("/v1/signin", requestBody))

	return res
}
