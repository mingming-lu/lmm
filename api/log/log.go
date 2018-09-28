package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultZapConfig returns zap config in default setting in this project
func DefaultZapConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg
}
