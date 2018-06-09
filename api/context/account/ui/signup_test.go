package ui

import (
	"io"
	"lmm/api/context/account/appservice"
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/repository"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestPostV1Signup(t *testing.T) {
	tester := testing.NewTester(t)

	router := testing.NewRouter()
	router.POST("/v1/signup", SignUp)

	name, password := uuid.New()[:31], uuid.New()
	reqeustBody := testing.StructToRequestBody(Auth{Name: name, Password: password})

	res := testing.NewResponse()
	router.ServeHTTP(res, testing.POST("/v1/signup", reqeustBody))

	tester.Is(http.StatusCreated, res.StatusCode())
	tester.Regexp(`/users/\d+`, res.Header().Get("Location"))
}

func TestPostV1Signup_Duplicate(t *testing.T) {
	tester := testing.NewTester(t)

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repository.New().Add(user)

	router := testing.NewRouter()
	router.POST("/v1/signup", SignUp)

	res := testing.NewResponse()
	auth := Auth{Name: name, Password: password}
	router.ServeHTTP(res, testing.POST("/v1/signup", testing.StructToRequestBody(auth)))

	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrDuplicateUserName.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyUserName(t *testing.T) {
	tester := testing.NewTester(t)

	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: "1234"})
	res := postSignUp(requestBody)

	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyPassword(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "foobar", Password: ""})
	res := postSignUp(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyUserNameAndPassword(t *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: ""})
	res := postSignUp(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(appservice.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func postSignUp(requestBody io.Reader) *testing.Response {
	res := testing.NewResponse()

	router := http.NewRouter()
	router.POST("/v1/signup", SignIn)
	router.ServeHTTP(res, testing.POST("/v1/signup", requestBody))

	return res
}
