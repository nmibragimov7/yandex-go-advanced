package middleware

import (
	"bytes"
	"fmt"
	"time"
	"yandex-go-advanced/internal/logger"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(data []byte) (int, error) {
	w.body.Write(data)

	n, err := w.ResponseWriter.Write(data)
	if err != nil {
		return n, fmt.Errorf("response writer error: %w", err)
	}
	return n, nil
}

func LoggerMiddleware(sgr *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method

		rbw := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = rbw

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()
		duration := time.Since(start)

		sugar := sgr.Get()
		sugar.Info(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", status,
			"size", size,
		)
	}
}
