package router

import (
	"time"
	"yandex-go-advanced/internal/common"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterProvider struct {
	Config  *config.Config
	Sugar   *zap.SugaredLogger
	Handler common.Handler
}

func (p *RouterProvider) Router() *gin.Engine {
	r := gin.Default()
	sugarWithCtx := p.Sugar.With(
		"app", "shortener",
		"service", "main",
		"func", "Router",
	)

	r.Use(middleware.GzipMiddleware(sugarWithCtx))
	r.Use(middleware.LoggerMiddleware(sugarWithCtx))
	r.Use(middleware.TimeoutMiddleware(sugarWithCtx, 2*time.Second))

	r.POST("/", p.Handler.MainPage)
	r.POST("/api/shorten", p.Handler.ShortenHandler)
	r.GET("/ping", p.Handler.PingHandler)
	r.POST("/api/shorten/batch", p.Handler.ShortenBatchHandler)
	r.GET("/:id", p.Handler.IDPage)

	return r
}
