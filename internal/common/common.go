package common

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GzipMiddleware interface {
	GzipMiddleware(sgr *zap.SugaredLogger) gin.HandlerFunc
}

type LoggerMiddleware interface {
	LoggerMiddleware(sgr *zap.SugaredLogger) gin.HandlerFunc
}

type Handler interface {
	MainPage(c *gin.Context)
	IDPage(c *gin.Context)
	ShortenHandler(c *gin.Context)
}
