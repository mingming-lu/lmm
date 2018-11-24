package middleware

import (
	"lmm/api/http"
	"lmm/api/testing"
)

func TestCacheControl(tt *testing.T) {
	t := testing.NewTester(tt)

	router := http.NewRouter()
	router.Use(CacheControl)

	router.GET("/", func(c http.Context) {
		c.String(http.StatusOK, "OK")
	})

	req := testing.GET("/", nil)
	res := testing.NewResponse()
	router.ServeHTTP(res, req)

	t.Is(res.StatusCode(), http.StatusOK)
	t.Is(res.Body(), "OK")
	t.Is(res.Header().Get("Cache-Control"), "public, max-age=60")
}
