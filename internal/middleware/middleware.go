package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GzipProvider struct{}
type LoggerProvider struct{}

type gzipWriter struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	n, err := w.zw.Write(b)
	if err != nil {
		return 0, fmt.Errorf("failed to write compressed data: %w", err)
	}
	return n, nil
}

func (w *gzipWriter) Close() error {
	err := w.zw.Close()
	if err != nil {
		return fmt.Errorf("failed to close compressed data: %w", err)
	}
	return nil
}

const (
	logKeyError = "error"
)

func (p *GzipProvider) GzipMiddleware(sgr *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
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
					sgr.Errorw(
						"gzip middleware write close failed",
						logKeyError, err.Error(),
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
		if sendsGzip {
			zr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				sgr.Errorw(
					"gzip middleware reader failed",
					logKeyError, err.Error(),
				)
				c.Writer.WriteHeader(http.StatusBadRequest)
				_, err = c.Writer.WriteString(http.StatusText(http.StatusBadRequest))
				if err != nil {
					sgr.Errorw(
						"gzip middleware write failed",
						logKeyError, err.Error(),
					)
				}
				c.Abort()
				return
			}
			defer func() {
				err := zr.Close()
				if err != nil {
					sgr.Errorw(
						"gzip middleware reader close failed",
						logKeyError, err.Error(),
					)
				}
			}()

			c.Request.Body = io.NopCloser(zr)
		}

		c.Next()
	}
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

func (p *LoggerProvider) LoggerMiddleware(sgr *zap.SugaredLogger) gin.HandlerFunc {
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

		sgr.Infow(
			"request handler log",
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", status,
			"size", size,
		)
	}
}
