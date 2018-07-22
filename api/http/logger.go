package http

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Logger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type DefaultLogger struct {
	mux         sync.Mutex
	callerDepth int

	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

func newLogger() *DefaultLogger {
	logger := &DefaultLogger{
		callerDepth: 2,

		info: log.New(os.Stdout, "[info] ", log.LstdFlags|log.Llongfile),
		warn: log.New(os.Stderr, "[warn] ", log.LstdFlags|log.Llongfile),
		err:  log.New(os.Stderr, "[error] ", log.LstdFlags|log.Llongfile),
	}
	return logger
}

func (l *DefaultLogger) SetCallerDepth(depth int) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.callerDepth = depth
}

func (l *DefaultLogger) Info(format string, args ...interface{}) {
	l.info.Output(l.callerDepth, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Warn(format string, args ...interface{}) {
	l.warn.Output(l.callerDepth, fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Error(format string, args ...interface{}) {
	l.err.Output(l.callerDepth, fmt.Sprintf(format, args...))
}
