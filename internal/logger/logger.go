package logger

import (
	"log"

	"go.uber.org/zap"
)

func Init() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Zap NewDevelopment: %s", err.Error())
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Printf("Logger Sync: %s", err.Error())
		}
	}()

	return logger.Sugar()
}
