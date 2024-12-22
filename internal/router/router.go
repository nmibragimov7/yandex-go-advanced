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
	str *storage.FileStorage,
	sgr *logger.Logger,
	mp common.Middleware,
	hp common.Handler,
) *gin.Engine {
	r := gin.Default()

	r.Use(mp.GzipMiddleware(sgr))
	r.Use(mp.LoggerMiddleware(sgr))
	r.POST("/", func(c *gin.Context) {
		hp.MainPage(c, cnf, str, sgr)
	})
	r.POST("/api/shorten", func(c *gin.Context) {
		hp.ShortenHandler(c, cnf, str, sgr)
	})
	r.GET("/:id", func(c *gin.Context) {
		hp.IDPage(c, str, sgr)
	})

	return r
}
