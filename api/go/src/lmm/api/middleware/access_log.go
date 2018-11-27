package middleware

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lmm/api/http"
)

type kafkaWriter struct{}

func newKafkaSyncWriter() zapcore.WriteSyncer {
	w := zapcore.AddSync(new(kafkaWriter))
	return zapcore.Lock(w)
}

func (w *kafkaWriter) Write(b []byte) (int, error) {
	zap.L().Info("TODO: send data to kafka",
		zap.ByteString("data", b),
	)

	return 0, nil
}

var (
	stdoutEnabler = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv == zapcore.InfoLevel
	})
	stderrEnabler = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.WarnLevel
	})
)

func newAccessLog() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:     "level",
		MessageKey:   "msg",
		NameKey:      "logger",
		TimeKey:      "ts",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	kafkaEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, newKafkaSyncWriter(), stderrEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), stderrEnabler),
		zapcore.NewCore(kafkaEncoder, newKafkaSyncWriter(), stdoutEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), stdoutEnabler),
	)
	return zap.New(core).Named("access_log")
}

// AccessLog records access log
func AccessLog(next http.Handler) http.Handler {
	logger := newAccessLog()
	defer logger.Sync()

	return func(c http.Context) {
		start := time.Now()

		next(c)

		req := c.Request()
		res := c.Response()
		status := res.StatusCode()
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("request_id", req.RequestID()),
			zap.String("remote_addr", req.RemoteAddr()),
			zap.String("ua", req.UserAgent()),
			zap.String("method", req.Method),
			zap.String("host", req.Host()),
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
