package ui

import (
	"log"

	"lmm/api/http"
	"lmm/api/service/account/domain/factory"
	"lmm/api/service/account/domain/model"
	"lmm/api/service/account/domain/service"
	"lmm/api/service/account/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestBearerAuth_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repo.Add(user)

	handler := ui.BearerAuth(func(c *http.Context) {
		field, ok := c.Values().Get("user").(*model.User)
		t.True(ok)
		t.Is(user.ID(), field.ID())
	})

	headers := make(map[string]string, 0)
	headers["Authorization"] = "Bearer " + service.EncodeToken(user.Token())
	res := request(headers, handler)
	t.Is(http.StatusOK, res.StatusCode())
}

func TestBearerAuth_NoAuthorization(t *testing.T) {
	tester := testing.NewTester(t)
	repo := infra.NewUserStorage(testing.DB())

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repo.Add(user)

	handler := ui.BearerAuth(func(c *http.Context) {
		log.Println("This message should not to be shown")
	})
	var res *testing.Response
	tester.Output("", func() {
		res = request(nil, handler)
	})
	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestBearerAuth_NotBearerAuthorization(tt *testing.T) {
	t := testing.NewTester(tt)

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)

	handler := ui.BearerAuth(func(c *http.Context) {
		log.Println("This message should not to be shown")
	})
	var res *testing.Response
	t.Output("", func() {
		headers := make(map[string]string, 0)
		headers["Authorization"] = "Basic " + service.EncodeToken(user.Token())
		res = request(headers, handler)
	})
	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestBearerAuth_InvaidToken(t *testing.T) {
	tester := testing.NewTester(t)
	handler := ui.BearerAuth(func(c *http.Context) {
		log.Println("This message should not to be shown")
	})
	var res *testing.Response
	tester.Output("", func() {
		headers := make(map[string]string, 0)
		headers["Authorization"] = "Bearer xxx"
		res = request(headers, handler)
	})
	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func request(headers map[string]string, handler http.Handler) *testing.Response {
	request := testing.GET("/dosomething")
	response := testing.NewResponse()
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.GET("/dosomething", handler)
	router.ServeHTTP(response, request)

	return response
}
