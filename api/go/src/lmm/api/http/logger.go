package http

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	cfg := DefaultZapConfig()
	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	ReplaceLogger(&loggerImpl{core: zapLogger.Named("http")})
}

// Logger defines interface for logging
type Logger interface {
	Info(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
	Panic(context.Context, string)
}

// DefaultZapConfig returns zap config in default setting in this project
func DefaultZapConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg
}

var mu sync.Mutex
var globalLogger Logger

// ReplaceLogger replaces http logger thread safely
func ReplaceLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

// Info log format
func Info(c context.Context, msg string) {
	globalLogger.Info(c, msg)
}

// Warn log format
func Warn(c context.Context, msg string) {
	globalLogger.Warn(c, msg)
}

// Error log format
func Error(c context.Context, msg string) {
	globalLogger.Warn(c, msg)
}

func Panic(c context.Context, msg string) {
	globalLogger.Panic(c, msg)
}

type loggerImpl struct {
	core *zap.Logger
}

func (l *loggerImpl) Info(c context.Context, msg string) {
	reqID := extractRequestID(c)
	l.core.With(zap.String("request_id", reqID)).Info(msg)

}

func (l *loggerImpl) Warn(c context.Context, msg string) {
	reqID := extractRequestID(c)
	l.core.With(zap.String("request_id", reqID)).Warn(msg)
}

func (l *loggerImpl) Error(c context.Context, msg string) {
	reqID := extractRequestID(c)
	l.core.With(zap.String("request_id", reqID)).Error(msg)
}

func (l *loggerImpl) Panic(c context.Context, msg string) {
	reqID := extractRequestID(c)
	l.core.With(zap.String("request_id", reqID)).Panic(msg)
}
