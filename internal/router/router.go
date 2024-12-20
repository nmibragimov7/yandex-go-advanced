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

	routes := r.Use(mp.GzipMiddleware(sgr))
	routes.Use(mp.LoggerMiddleware(sgr))
	routes.POST("/", func(c *gin.Context) {
		hp.MainPage(c, cnf, str, sgr)
	})
	routes.GET("/:id", func(c *gin.Context) {
		hp.IDPage(c, str, sgr)
	})
	routes.POST("/api/shorten", func(c *gin.Context) {
		hp.ShortenHandler(c, cnf, str, sgr)
	})

	return r
}
