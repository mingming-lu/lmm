package middleware

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"lmm/api/http"
)

// Recovery tries to recover panics
func Recovery(next http.Handler) http.Handler {
	return func(c http.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if s, ok := recovered.(string); ok {

					// occurred when client abort
					if strings.Contains(s, context.Canceled.Error()) {
						http.ClientAbort(c)
						return
					}

				}

				fields := []zap.Field{
					zap.String("request_id", c.Request().RequestID()),
					zap.Any("what", recovered),
				}
				zap.L().Error("unexpected error", fields...)
				http.InternalServerError(c)
			}
		}()
		next(c)
	}
}
