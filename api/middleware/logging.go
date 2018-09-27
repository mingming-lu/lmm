package middleware

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lmm/api/http"
)

var logger *zap.Logger

func init() {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if l, err := cfg.Build(); err == nil {
		logger = l
	} else {
		panic(err)
	}
}

// Logging middleware logs
func Logging(next http.Handler) http.Handler {
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
			zap.String("path", req.RequestURI),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("ua", req.Header.Get("User-Agent")),
			zap.String("latency", time.Since(start).String()),
		}

		if status >= 500 {
			logger.Error(http.StatusText(status), fields...)
		} else if status >= 400 {
			logger.Warn(http.StatusText(status), fields...)
		} else {
			logger.Info(http.StatusText(status), fields...)
		}
	}
}
