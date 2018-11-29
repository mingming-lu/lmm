package middleware

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lmm/api/http"
)

var (
	stdoutEnabler = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv == zapcore.InfoLevel
	})
	stderrEnabler = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.WarnLevel
	})
)

// AccessLogger struct
type AccessLogger struct {
	logger *zap.Logger
	writer io.Writer
}

// NewAccessLog creates a new AccessLog
func NewAccessLog(logWriter io.Writer) *AccessLogger {
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:     "level",
		MessageKey:   "msg",
		NameKey:      "logger",
		TimeKey:      "ts",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		LineEnding:   zapcore.DefaultLineEnding,
	}
	writerEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(writerEncoder, newSyncWriter(logWriter), stderrEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), stderrEnabler),
		zapcore.NewCore(writerEncoder, newSyncWriter(logWriter), stdoutEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), stdoutEnabler),
	)
	core = zapcore.NewSampler(core, time.Second, 100, 100)

	return &AccessLogger{
		logger: zap.New(core).Named("access_log"),
		writer: logWriter,
	}
}

// Sync implementation
func (l *AccessLogger) Sync() error {
	return l.logger.Sync()
}

// Write may be not mutex
func (l *AccessLogger) Write(p []byte) (int, error) {
	return l.Write(p)
}

// AccessLog records access log
func (l *AccessLogger) AccessLog(next http.Handler) http.Handler {
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
			l.logger.Error(http.StatusText(status), fields...)
		} else if status >= 400 {
			l.logger.Warn(http.StatusText(status), fields...)
		} else {
			l.logger.Info(http.StatusText(status), fields...)
		}
	}
}

func newSyncWriter(w io.Writer) zapcore.WriteSyncer {
	syncWriter := zapcore.AddSync(w)
	return zapcore.Lock(syncWriter)
}
