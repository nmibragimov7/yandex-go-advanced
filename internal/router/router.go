package router

import (
	"time"
	"yandex-go-advanced/internal/common"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/middleware"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterProvider struct {
	Storage storage.Storage
	Config  *config.Config
	Sugar   *zap.SugaredLogger
	Handler common.Handler
	Session *session.SessionProvider
}

func (p *RouterProvider) Router() *gin.Engine {
	r := gin.Default()
	sugarWithCtx := p.Sugar.With(
		"app", "shortener",
		"service", "main",
		"func", "Router",
	)

	middlewares := []gin.HandlerFunc{
		middleware.GzipMiddleware(sugarWithCtx),
		middleware.LoggerMiddleware(sugarWithCtx),
		middleware.TimeoutMiddleware(sugarWithCtx, 2*time.Second),
	}

	r.Use(middlewares...)
	if *p.Config.DataBase != "" {
		r.Use(middleware.AuthMiddleware(sugarWithCtx, p.Storage, p.Session))
	}

	r.POST("/", p.Handler.MainPage)
	r.POST("/api/shorten", p.Handler.ShortenHandler)
	r.GET("/ping", p.Handler.PingHandler)
	r.POST("/api/shorten/batch", p.Handler.ShortenBatchHandler)
	r.GET("/api/user/urls", p.Handler.UserUrlsHandler)
	r.GET("/:id", p.Handler.IDPage)

	return r
}
