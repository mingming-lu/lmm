package middleware

import (
	"lmm/api/http"
)

// CacheControl adds Cache-Control header to responses
func CacheControl(next http.Handler) http.Handler {
	return func(c http.Context) {
		c.Response().Header().Add("Cache-Control", "public, max-age=60")
		next(c)
	}
}
