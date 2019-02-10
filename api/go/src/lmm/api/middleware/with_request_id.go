package middleware

import (
	"lmm/api/http"
	"lmm/api/util/contextutil"
)

// WithRequestID adds request id value to context
func WithRequestID(next http.Handler) http.Handler {
	return func(c http.Context) {
		requestID := c.Request().RequestID()
		if requestID == "" {
			http.Log().Warn(c, "empty request id !")
			http.BadRequest(c)
			return
		}

		ctx := contextutil.WithRequestID(c.Request().Context(), requestID)

		next(c.With(ctx))
	}
}
