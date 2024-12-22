package common

import (
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	GzipMiddleware(sgr *logger.Logger) gin.HandlerFunc
	LoggerMiddleware(sgr *logger.Logger) gin.HandlerFunc
}

type Handler interface {
	MainPage(c *gin.Context, cnf *config.Config, str *storage.FileStorage, sgr *logger.Logger)
	IDPage(c *gin.Context, str *storage.FileStorage, sgr *logger.Logger)
	ShortenHandler(c *gin.Context, cnf *config.Config, str *storage.FileStorage, sgr *logger.Logger)
}
