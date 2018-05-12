package user

import (
	"lmm/api/testing"
	"lmm/api/usecase/user"
	"net/http"
)

func TestPostV1Signup(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	router := testing.NewRouter()
	router.POST("/v1/signup", SignUp)

	reqeustBody := testing.StructToRequestBody(user.Auth{Name: "foobar", Password: "1234"})

	res := testing.NewResponse()
	router.ServeHTTP(res, testing.POST("/v1/signup", reqeustBody))

	tester.Is(res.StatusCode(), http.StatusCreated)
	tester.Is(res.Header().Get("Location"), "/users/1")
}
