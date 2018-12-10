package middleware

import (
	"context"

	"lmm/api/http"
)

// WithRequestID adds request id value to context
func WithRequestID(next http.Handler) http.Handler {
	return func(c http.Context) {
		requestID := c.Request().RequestID()
		if requestID == "" {
			http.Warn(c, "empty request id !")
			http.BadRequest(c)
			return
		}

		ctx := context.WithValue(c.Request().Context(), http.RequestIDContextKey, requestID)

		next(c.With(ctx))
	}
}
