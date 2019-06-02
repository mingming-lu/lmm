package http

import (
	"context"

	"google.golang.org/appengine/log"
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
	log.Infof(c, msg)
}

func (l *loggerImpl) Warn(c context.Context, msg string) {
	log.Warningf(c, msg)
}

func (l *loggerImpl) Error(c context.Context, msg string) {
	log.Errorf(c, msg)
}

func (l *loggerImpl) Panic(c context.Context, msg string) {
	log.Criticalf(c, msg)
}
