package router

import (
	"net/http"
	"time"

	"yandex-go-advanced/internal/common"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/middleware"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RouterProvider - struct that contains the necessary router settings
type RouterProvider struct {
	Storage storage.Storage
	Config  *config.Config
	Sugar   *zap.SugaredLogger
	Handler common.Handler
	Session *session.SessionProvider
}

// Router - func for init router
func (p *RouterProvider) Router() *gin.Engine {
	r := gin.Default()
	sugarWithCtx := p.Sugar.With(
		"app", "shortener",
		"service", "main",
		"func", "Router",
	)
	r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	middlewares := []gin.HandlerFunc{
		middleware.GzipMiddleware(sugarWithCtx),
		middleware.LoggerMiddleware(sugarWithCtx),
		middleware.TimeoutMiddleware(sugarWithCtx, 2*time.Second),
	}

	r.Use(middlewares...)

	atp := &middleware.AuthProvider{
		Storage: p.Storage,
		Config:  p.Config,
		Sugar:   p.Sugar,
		Session: p.Session,
	}
	r.POST("/", middleware.AuthMiddleware(atp), p.Handler.MainPage)
	r.POST("/api/shorten", middleware.AuthMiddleware(atp), p.Handler.ShortenHandler)
	r.POST("/api/shorten/batch", middleware.AuthMiddleware(atp), p.Handler.ShortenBatchHandler)
	r.GET("/api/user/urls", middleware.AuthMiddleware(atp), p.Handler.UserUrlsHandler)
	r.DELETE("/api/user/urls", p.Handler.UserUrlsDeleteHandler)
	r.GET("/ping", p.Handler.PingHandler)
	r.GET("/:id", p.Handler.IDPage)
	r.GET("/api/internal/stats", p.Handler.TrustedSubnetHandler)

	return r
}
