package middleware

import (
	"log"
	"runtime/debug"

	"lmm/api/http"
)

// Recovery is a middleware which handles panic event
func Recovery(next http.Handler) http.Handler {
	return func(c http.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// TODO specify stack tracing depth
				log.Print(recovered)
				debug.PrintStack()
				http.InternalServerError(c)
			}
		}()
		next(c)
	}
}
