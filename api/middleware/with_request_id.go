package middleware

import (
	"context"

	"lmm/api/http"
)

// WithRequestID adds request id value to  context
func WithRequestID(next http.Handler) http.Handler {
	return func(c http.Context) {
		requestIDKey := http.StrCtxKey("request_id")
		requestID := c.Request().RequestID()
		if requestID == "" {
			requestID = "-"
		}
		ctx := context.WithValue(c.Request().Context(), requestIDKey, requestID)
		next(c.With(ctx))
	}
}
