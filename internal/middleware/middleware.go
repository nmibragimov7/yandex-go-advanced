package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"yandex-go-advanced/internal/logger"

	"github.com/gin-gonic/gin"
)

type Provider struct{}

type gzipWriter struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.zw.Write(b)
}

func (w *gzipWriter) Close() error {
	return w.zw.Close()
}

type gzipReader struct {
	io.Reader
	zr *gzip.Reader
}

func (r *gzipReader) Read(b []byte) (int, error) {
	return r.zr.Read(b)
}

func (r *gzipReader) Close() error {
	return r.zr.Close()
}

type loggerWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *loggerWriter) Write(data []byte) (int, error) {
	w.body.Write(data)

	n, err := w.ResponseWriter.Write(data)
	if err != nil {
		return n, fmt.Errorf("response writer error: %w", err)
	}
	return n, nil
}

func (p *Provider) GzipMiddleware(sgr *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sugar := sgr.Get()

		contentType := c.Request.Header.Get("Content-Type")
		supportsJSON := strings.Contains(contentType, "application/json")
		supportsHTML := strings.Contains(contentType, "text/html")

		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip && (supportsJSON || supportsHTML) {
			zw := gzip.NewWriter(c.Writer)
			defer func() {
				err := zw.Close()
				if err != nil {
					sugar.Errorw(
						"gzip middleware write close failed",
						"error", err.Error(),
					)
				}
			}()
			c.Writer = &gzipWriter{
				ResponseWriter: c.Writer,
				zw:             zw,
			}
			c.Writer.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip && (supportsJSON || supportsHTML) {
			zr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				sugar.Errorw(
					"gzip middleware reader failed",
					"error", err.Error(),
				)
				c.Writer.WriteHeader(http.StatusBadRequest)
				_, err = c.Writer.WriteString(http.StatusText(http.StatusBadRequest))
				if err != nil {
					sugar.Errorw(
						"gzip middleware write failed",
						"error", err.Error(),
					)
				}
				c.Abort()
				return
			}
			defer func() {
				err := zr.Close()
				if err != nil {
					sugar.Errorw(
						"gzip middleware reader close failed",
						"error", err.Error(),
					)
				}
			}()

			c.Request.Body = &gzipReader{
				Reader: zr,
				zr:     zr,
			}
		}

		c.Next()
	}
}

func (p *Provider) LoggerMiddleware(sgr *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method

		rbw := &loggerWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = rbw

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()
		duration := time.Since(start)

		sugar := sgr.Get()
		sugar.Infow(
			"request handler log",
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", status,
			"size", size,
		)
	}
}
