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
	MainPage(c *gin.Context, cnf *config.Config, str *storage.Store, sgr *logger.Logger)
	IDPage(c *gin.Context, str *storage.Store, sgr *logger.Logger)
	ShortenHandler(c *gin.Context, cnf *config.Config, str *storage.Store, sgr *logger.Logger)
}
