package router

import (
	"yandex-go-advanced/internal/common"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/middleware"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterProvider struct {
	Config  *config.Config
	Storage *storage.FileStorage
	Sugar   *zap.SugaredLogger
	Handler common.Handler
}

func (p *RouterProvider) Router() *gin.Engine {
	r := gin.Default()
	sugarWithCtx := p.Sugar.With(
		"app", "shortener",
		"service", "main",
	)

	r.Use(middleware.GzipMiddleware(sugarWithCtx))
	r.Use(middleware.LoggerMiddleware(sugarWithCtx))

	r.POST("/", p.Handler.MainPage)
	r.POST("/api/shorten", p.Handler.ShortenHandler)
	r.GET("/ping", p.Handler.PingHandler)
	r.GET("/:id", p.Handler.IDPage)

	return r
}
