package router

import (
	"yandex-go-advanced/internal/common"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Provider struct {
	Config           *config.Config
	Storage          *storage.FileStorage
	Sugar            *zap.SugaredLogger
	GzipMiddleware   common.GzipMiddleware
	LoggerMiddleWare common.LoggerMiddleware
	Handler          common.Handler
}

func (p *Provider) Router() *gin.Engine {
	r := gin.Default()
	sugarWithCtx := p.Sugar.With(
		"app", "shortener",
		"service", "main",
	)

	r.Use(p.GzipMiddleware.GzipMiddleware(sugarWithCtx))
	r.Use(p.LoggerMiddleWare.LoggerMiddleware(sugarWithCtx))

	r.POST("/", p.Handler.MainPage)
	r.POST("/api/shorten", p.Handler.ShortenHandler)
	r.GET("/:id", p.Handler.IDPage)

	return r
}
