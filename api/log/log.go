package log

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var defaultLogger *zap.Logger

func newDefaultLogger() {
	if env := os.Getenv("ENV"); env == "PRODUCTION" {
		if logger, err := zap.NewProduction(); err == nil {
			defaultLogger = logger
		} else {
			log.Fatal(err)
		}
	} else {
		if logger, err := zap.NewDevelopment(); err != nil {
			defaultLogger = logger
		} else {
			log.Fatal(err)
		}
	}
}
