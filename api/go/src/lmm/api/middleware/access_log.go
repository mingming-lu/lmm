package middleware

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lmm/api/http"
)

// AccessLog records access log
func AccessLog(next http.Handler) http.Handler {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = logger.Named("access_log")

	return func(c http.Context) {
		start := time.Now()

		next(c)

		req := c.Request()
		res := c.Response()
		status := res.StatusCode()
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("request_id", req.RequestID()),
			zap.String("client_ip", req.ClientIP()),
			zap.String("forwarded_for", req.Header.Get("X-Forwarded-For")),
			zap.String("ua", req.UserAgent()),
			zap.String("method", req.Method),
			zap.String("host", req.HostName()),
			zap.String("uri", req.RequestURI),
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
