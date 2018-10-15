package middleware

import (
	"time"

	"go.uber.org/zap"

	"lmm/api/http"
)

// AccessLog records access log
func AccessLog(next http.Handler) http.Handler {
	return func(c http.Context) {
		start := time.Now()

		next(c)

		req := c.Request()
		res := c.Response()
		status := res.StatusCode()
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("request_id", req.RequestID()),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("ua", req.UserAgent()),
			zap.String("method", req.Method),
			zap.String("host", req.Host),
			zap.String("uri", req.RequestURI),
			zap.String("latency", time.Since(start).String()),
		}

		if status >= 500 {
			zap.L().Error(http.StatusText(status), fields...)
		} else if status >= 400 {
			zap.L().Warn(http.StatusText(status), fields...)
		} else {
			zap.L().Info(http.StatusText(status), fields...)
		}
	}
}
