package middleware

import (
	"os"

	"lmm/api/http"
)

var (
	appHost = os.Getenv("APP_HOST")
)

// CacheControl adds Cache-Control header to responses
func CacheControl(next http.Handler) http.Handler {
	return func(c http.Context) {
		if c.Request().Origin() == appHost && c.Request().Method == "GET" {
			c.Response().Header().Add("Cache-Control", "public, max-age=60")
		}
		next(c)
	}
}
