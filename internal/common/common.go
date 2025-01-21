package common

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	MainPage(c *gin.Context)
	IDPage(c *gin.Context)
	ShortenHandler(c *gin.Context)
	PingHandler(c *gin.Context)
	ShortenBatchHandler(c *gin.Context)
	UserUrlsHandler(c *gin.Context)
}
