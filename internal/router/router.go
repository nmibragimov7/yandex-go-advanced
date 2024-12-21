package router

import (
	"yandex-go-advanced/internal/common"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
)

func Router(
	cnf *config.Config,
	str *storage.Store,
	sgr *logger.Logger,
	mp common.Middleware,
	hp common.Handler,
) *gin.Engine {
	r := gin.Default()

	gzipLoggerGroup := r.Group("/")

	gzipLoggerGroup.Use(mp.GzipMiddleware(sgr))
	gzipLoggerGroup.Use(mp.LoggerMiddleware(sgr))
	gzipLoggerGroup.POST("/", func(c *gin.Context) {
		hp.MainPage(c, cnf, str, sgr)
	})
	gzipLoggerGroup.POST("/api/shorten", func(c *gin.Context) {
		hp.ShortenHandler(c, cnf, str, sgr)
	})

	loggerGroup := r.Group("/")
	loggerGroup.Use(mp.LoggerMiddleware(sgr))
	loggerGroup.GET("/:id", func(c *gin.Context) {
		hp.IDPage(c, str, sgr)
	})

	return r
}
