package middleware

import (
	"time"

	"go.uber.org/zap"

	"lmm/api/http"
)

// NewAccessLog returns a middleware to record access log
func NewAccessLog(logger *zap.Logger) http.Middleware {
	r := accessLogRecorder{logger: logger}
	return r.accessLog
}

type accessLogRecorder struct {
	logger *zap.Logger
}

func (r *accessLogRecorder) accessLog(next http.Handler) http.Handler {
	return func(c http.Context) {
		start := time.Now()

		next(c)

		req := c.Request()
		res := c.Response()

		status := res.StatusCode()
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("request_id", req.Header.Get("X-Request-ID")),
			zap.String("method", req.Method),
			zap.String("proto", req.Proto),
			zap.String("host", req.Host),
			zap.String("uri", req.RequestURI),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("ua", req.Header.Get("User-Agent")),
			zap.String("latency", time.Since(start).String()),
		}

		if status >= 500 {
			r.logger.Error(http.StatusText(status), fields...)
		} else if status >= 400 {
			r.logger.Warn(http.StatusText(status), fields...)
		} else {
			r.logger.Info(http.StatusText(status), fields...)
		}
	}
}
