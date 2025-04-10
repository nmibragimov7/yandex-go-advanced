package logger

import (
	"log"

	"go.uber.org/zap"
)

var newProduction = zap.NewProduction

// Init - initialize logger instance
func Init() *zap.SugaredLogger {
	logger, err := newProduction()
	if err != nil {
		log.Printf("failed to build development logger: %s", err.Error())
		return nil
	}

	return logger.Sugar()
}
