package http

import (
	"context"

	"lmm/api/util/contextutil"

	"go.uber.org/zap"
)

var logger *loggerImpl

func init() {
	logger = new(loggerImpl)
}

// Logger defines interface for logging
type Logger interface {
	Info(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
	Panic(context.Context, string)
}

// Log gets the singleton of default logger implementation
func Log() Logger {
	return logger
}

type loggerImpl struct{}

func (l *loggerImpl) Info(c context.Context, msg string) {
	zap.L().Info(msg, zap.String("request_id", contextutil.RequestID(c)))
}

func (l *loggerImpl) Warn(c context.Context, msg string) {
	zap.L().Warn(msg, zap.String("request_id", contextutil.RequestID(c)))
}

func (l *loggerImpl) Error(c context.Context, msg string) {
	zap.L().Error(msg, zap.String("request_id", contextutil.RequestID(c)))
}

func (l *loggerImpl) Panic(c context.Context, msg string) {
	zap.L().Panic(msg, zap.String("request_id", contextutil.RequestID(c)))
}
