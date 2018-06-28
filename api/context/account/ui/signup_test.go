package ui

import (
	"io"
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/service"
	"lmm/api/context/account/infra"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestPostV1Signup(t *testing.T) {
	tester := testing.NewTester(t)

	auth := Auth{
		Name:     uuid.New()[:31],
		Password: uuid.New(),
	}

	res := postSignUp(testing.StructToRequestBody(auth))

	tester.Is(http.StatusCreated, res.StatusCode())
	tester.Regexp(`/users/\d+`, res.Header().Get("Location"))
}

func TestPostV1Signup_Duplicate(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)

	t.NoError(repo.Add(user))

	router := testing.NewRouter()
	router.POST("/v1/signup", ui.SignUp)

	res := testing.NewResponse()
	auth := Auth{Name: name, Password: password}
	router.ServeHTTP(res, testing.POST("/v1/signup", testing.StructToRequestBody(auth)))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(service.ErrDuplicateUserName.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyUserName(tt *testing.T) {
	t := testing.NewTester(tt)

	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: "1234"})
	res := postSignUp(requestBody)

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyPassword(tt *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "foobar", Password: ""})
	res := postSignUp(requestBody)

	t := testing.NewTester(tt)
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyUserNameAndPassword(tt *testing.T) {
	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: ""})
	res := postSignUp(requestBody)

	t := testing.NewTester(tt)
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(service.ErrInvalidUserNameOrPassword.Error()+"\n", res.Body())
}

func postSignUp(requestBody io.ReadCloser) *testing.Response {
	res := testing.NewResponse()

	router := http.NewRouter()
	router.POST("/v1/signup", ui.SignUp)
	router.ServeHTTP(res, testing.POST("/v1/signup", requestBody))

	return res
}
