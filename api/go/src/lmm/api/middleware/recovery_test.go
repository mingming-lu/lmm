package middleware

import (
	"context"
	"lmm/api/http"
	"lmm/api/testing"
)

func TestRecovery(tt *testing.T) {
	t := testing.NewTester(tt)

	router := http.NewRouter()
	router.Use(Recovery)
	router.GET("/", func(c http.Context) {
		if xPanic := c.Request().Header.Get("X-Panic"); xPanic != "" {
			http.Log().Panic(c, xPanic)
		}
		c.String(http.StatusOK, "no panic")
	})

	t.Run("NormalPanic", func(_ *testing.T) {
		req := testing.GET("/", nil)
		req.Header.Add("X-Panic", "panic!!!")
		res := testing.DoRequest(req, router)

		t.Is(http.StatusInternalServerError, res.StatusCode())
		t.Is(http.StatusText(http.StatusInternalServerError), res.Body())
	})

	t.Run("ContextCanceled", func(_ *testing.T) {
		req := testing.GET("/", nil)
		req.Header.Add("X-Panic", context.Canceled.Error())
		res := testing.DoRequest(req, router)

		t.Is(http.StatusClientAbort, res.StatusCode())
		t.Is(http.StatusText(http.StatusClientAbort), res.Body())
	})

	t.Run("NoPanic", func(_ *testing.T) {
		req := testing.GET("/", nil)
		res := testing.DoRequest(req, router)

		t.Is(http.StatusOK, res.StatusCode())
		t.Is("no panic", res.Body())
	})

}
