package log

import (
	"io"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init initializes logger for std log and zap global logger
// Expected to be called in main() first
func Init(logWriter io.Writer) func() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		NameKey:       "logger",
		TimeKey:       "ts",
		StacktraceKey: "trace",
		CallerKey:     "caller",
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		LineEnding:    zapcore.DefaultLineEnding,
	}

	kafkaEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	globalEnabler := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return true
	})

	syncer := zapcore.Lock(zapcore.AddSync(logWriter))

	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, syncer, globalEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), globalEnabler),
	)
	core = zapcore.NewSampler(core, time.Second, 100, 100)

	logger := zap.New(core).
		Named("logger").
		WithOptions(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
				return lv == zap.PanicLevel
			})),
		)

	undo := zap.ReplaceGlobals(logger)

	return func() {
		logger.Sync()
		undo()
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}
