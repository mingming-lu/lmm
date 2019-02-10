package middleware

import (
	"io"
	"os"
	"sync"
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
	name := os.Getenv("LOGGER_NAME_API_ACCESS_LOG")
	if name == "" {
		panic("empty access logger name")
	}

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

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, newSyncWriter(logWriter), stderrEnabler),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), stderrEnabler),
		zapcore.NewCore(encoder, newSyncWriter(logWriter), stdoutEnabler),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), stdoutEnabler),
	)
	core = zapcore.NewSampler(core, time.Second, 100, 100)

	return &AccessLogger{
		logger: zap.New(core).Named(name),
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

type accessLogFields struct {
	fields []zap.Field
}

func newAccessLogFields() *accessLogFields {
	f := accessLogFields{
		fields: make([]zap.Field, 9, 9),
	}
	return &f
}

var zapFieldsPool = sync.Pool{
	New: func() interface{} {
		return newAccessLogFields()
	},
}

func buildAccessLogFields(status int, reqID, ip, xff, ua, method, host, uri, latency string) *accessLogFields {
	f, ok := zapFieldsPool.Get().(*accessLogFields)
	if !ok {
		panic("expected a *accessLogFields")
	}

	f.fields[0].Type = zapcore.Int64Type
	f.fields[0].Key = "status"
	f.fields[0].Integer = int64(status)

	f.fields[1].Type = zapcore.StringType
	f.fields[1].Key = "request_id"
	f.fields[1].String = reqID

	f.fields[2].Type = zapcore.StringType
	f.fields[2].Key = "client_ip"
	f.fields[2].String = ip

	f.fields[3].Type = zapcore.StringType
	f.fields[3].Key = "forwarded_for"
	f.fields[3].String = xff

	f.fields[4].Type = zapcore.StringType
	f.fields[4].Key = "ua"
	f.fields[4].String = ua

	f.fields[5].Type = zapcore.StringType
	f.fields[5].Key = "method"
	f.fields[5].String = method

	f.fields[6].Type = zapcore.StringType
	f.fields[6].Key = "host"
	f.fields[6].String = host

	f.fields[7].Type = zapcore.StringType
	f.fields[7].Key = "uri"
	f.fields[7].String = uri

	f.fields[8].Type = zapcore.StringType
	f.fields[8].Key = "latency"
	f.fields[8].String = latency

	return f
}

// AccessLog records access log
func (l *AccessLogger) AccessLog(next http.Handler) http.Handler {
	return func(c http.Context) {
		start := time.Now()

		next(c)

		req := c.Request()
		res := c.Response()
		status := res.StatusCode()
		f := buildAccessLogFields(
			status,
			req.RequestID(),
			req.ClientIP(),
			req.Header.Get("X-Forwarded-For"),
			req.UserAgent(),
			req.Method,
			req.HostName(),
			req.RequestURI,
			time.Since(start).String(),
		)

		if status >= 500 {
			l.logger.Error(http.StatusText(status), f.fields...)
		} else if status >= 400 {
			l.logger.Warn(http.StatusText(status), f.fields...)
		} else {
			l.logger.Info(http.StatusText(status), f.fields...)
		}

		zapFieldsPool.Put(f)
	}
}

func newSyncWriter(w io.Writer) zapcore.WriteSyncer {
	syncWriter := zapcore.AddSync(w)
	return zapcore.Lock(syncWriter)
}
