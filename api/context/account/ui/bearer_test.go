package auth

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/service"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func TestBearerAuth_Success(t *testing.T) {
	tester := testing.NewTester(t)
	user := testingService.NewUser()
	handler := BearerAuth(func(c *elesion.Context) {
		field, ok := c.Fields().Get("user").(*model.User)
		tester.True(ok)
		tester.Is(user.ID, field.ID)
	})
	headers := make(map[string]string, 0)
	headers["Authorization"] = "Bearer " + service.EncodeToken(user.Token)
	res := request(headers, handler)
	tester.Is(http.StatusOK, res.StatusCode())
}

func TestBearerAuth_NoAuthorization(t *testing.T) {
	tester := testing.NewTester(t)
	handler := BearerAuth(func(c *elesion.Context) {
		log.Println("This message should not to be shown")
	})
	var res *testing.Response
	tester.Output("", func() {
		res = request(nil, handler)
	})
	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestBearerAuth_NotBearerAuthorization(t *testing.T) {
	tester := testing.NewTester(t)
	user := testingService.NewUser()
	handler := BearerAuth(func(c *elesion.Context) {
		log.Println("This message should not to be shown")
	})
	var res *testing.Response
	tester.Output("", func() {
		headers := make(map[string]string, 0)
		headers["Authorization"] = "Basic " + service.EncodeToken(user.Token)
		res = request(headers, handler)
	})
	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestBearerAuth_InvaidToken(t *testing.T) {
	tester := testing.NewTester(t)
	handler := BearerAuth(func(c *elesion.Context) {
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

func request(headers map[string]string, handler elesion.Handler) *testing.Response {
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
