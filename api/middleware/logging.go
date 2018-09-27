package middleware

import (
	"time"

	"go.uber.org/zap"

	"lmm/api/http"
)

var logger *zap.Logger

func init() {
	if l, err := zap.NewProduction(); err == nil {
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

		logger.Info(http.StatusText(res.StatusCode()),
			zap.Int("status", res.StatusCode()),
			zap.String("latency", time.Since(start).String()),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
			zap.String("host", req.Host),
		)
	}
}
