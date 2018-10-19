package main

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/labstack/echo"
)

func accessLog(next echo.HandlerFunc) echo.HandlerFunc {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		req := c.Request()
		res := c.Response()
		status := res.Status // why always 200 ?
		if e, ok := err.(*echo.HTTPError); ok {
			status = e.Code
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("request_id", req.Header.Get("X-Request-ID")),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("ua", req.UserAgent()),
			zap.String("method", req.Method),
			zap.String("host", req.Host),
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
		return err
	}
}

func main() {
	e := echo.New()

	e.Use(accessLog)

	e.Static("/images", "/static/images")
	e.Static("/photos", "/static/images")

	e.Logger.Fatal(e.Start(":8003"))
}
