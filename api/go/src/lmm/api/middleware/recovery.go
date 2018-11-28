package middleware

import (
	"lmm/api/http"
)

// Recovery tries to recover panics
func Recovery(next http.Handler) http.Handler {
	return func(c http.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				http.InternalServerError(c)
			}
		}()
		next(c)
	}
}
