package middleware

import (
	"go.uber.org/zap"

	"lmm/api/http"
)

// NewRecovery returns a middleware to recover panic
func NewRecovery(logger *zap.Logger) http.Middleware {
	r := recoveryRecoder{logger: logger}
	return r.Recovery
}

type recoveryRecoder struct {
	logger *zap.Logger
}

func (r *recoveryRecoder) Recovery(next http.Handler) http.Handler {
	return func(c http.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				fields := []zap.Field{
					zap.String("request_id", c.Request().RequestID()),
					zap.Reflect("what", recovered),
				}
				r.logger.Error("unexpected error", fields...)
				http.InternalServerError(c)
			}
		}()
		next(c)
	}
}
