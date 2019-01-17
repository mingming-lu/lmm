package http

import (
	"lmm/api/testing"
	"net/http"
	"time"
)

func TestRouterNotFound(tt *testing.T) {
	t := testing.NewTester(tt)

	t.Run("DefaultNotFoundHandler", func(_ *testing.T) {
		router := NewRouter()

		req := testing.GET("/", nil)
		res := testing.DoRequest(req, router)

		t.Is(StatusNotFound, res.StatusCode())
		t.Is(StatusText(StatusNotFound), res.Body())
	})

	t.Run("SpecifiedNotFoundHandler", func(_ *testing.T) {
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
