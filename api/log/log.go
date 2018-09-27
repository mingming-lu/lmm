package log

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	"lmm/api/http"
)

var defaultLogger *zap.Logger

func init() {
	if env := os.Getenv("ENV"); env == "PRODUCTION" {
		if logger, err := zap.NewProduction(); err == nil {
			defaultLogger = logger
		} else {
			log.Fatal(err)
		}
	} else {
		if logger, err := zap.NewDevelopment(); err == nil {
			defaultLogger = logger
		} else {
			log.Fatal(err)
		}
	}
}

// LoggingService is a logging middleware
func LoggingService(next http.Handler) http.Handler {
	return func(c http.Context) {
		start := time.Now()
		next(c)

		req := c.Request()
		res := c.Response()

		defaultLogger.Info(http.StatusText(res.StatusCode()),
			zap.Int("status", res.StatusCode()),
			zap.String("latency", time.Since(start).String()),
			zap.String("method", req.Method),
			zap.String("uri", req.RequestURI),
			zap.String("host", req.Host),
		)
	}
}
