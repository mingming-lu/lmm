package http

import (
	"net/http"
	"time"

	"lmm/api/testing"
)

func TestRouterNotFound(tt *testing.T) {
	tt.Run("DefaultNotFoundHandler", func(tt *testing.T) {
		t := testing.NewTester(tt)
		router := NewRouter()

		req := testing.GET("/", nil)
		res := testing.DoRequest(req, router)

		t.Is(StatusNotFound, res.StatusCode())
		t.Is(StatusText(StatusNotFound), res.Body())
	})

	tt.Run("SpecifiedNotFoundHandler", func(tt *testing.T) {
		t := testing.NewTester(tt)
		plainText := "pretended OK"

		router := NewRouter()
		router.NotFound(func(c Context) {
			c.String(StatusOK, plainText)
		})

		req := testing.GET("/", nil)
		res := testing.DoRequest(req, router)

		t.Is(StatusOK, res.StatusCode())
		t.Is(plainText, res.Body())
	})
}

func TestHandleTimeout(tt *testing.T) {
	t := testing.NewTester(tt)

	router := NewRouter()
	router.GET("/timeout", func(c Context) {
		time.Sleep(5 * time.Second)
	})

	res := testing.DoRequest(testing.GET("/timeout", nil), router)
	t.Is(http.StatusRequestTimeout, res.StatusCode())
}
