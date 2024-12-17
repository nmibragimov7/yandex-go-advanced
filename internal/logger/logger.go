package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

func (l *Logger) Get() *zap.SugaredLogger {
	return l.Sugar
}

func InitLogger() *Logger {
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

	return &Logger{
		Sugar: logger.Sugar(),
	}
}
