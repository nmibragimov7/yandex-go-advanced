package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	logKeyError = "error"
	cookieName  = "user_token"
)

func AuthMiddleware(sgr *zap.SugaredLogger, str storage.Storage, ssn *session.SessionProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(cookieName)
		if err != nil || cookie == "" {
			record := &models.UserRecord{}

			id, err := str.Set("users", record)
			if err != nil {
				sgr.Errorw(
					"failed to save user record",
					logKeyError, err.Error(),
				)

				message := models.Response{
					Message: http.StatusText(http.StatusInternalServerError),
				}

				c.JSON(http.StatusInternalServerError, message)
				c.Abort()
				return
			}

			token, err := ssn.GenerateToken(id.(int64))
			if err != nil {
				sgr.Errorw(
					"failed to generate token",
					logKeyError, err.Error(),
				)

				message := models.Response{
					Message: http.StatusText(http.StatusInternalServerError),
				}

				c.JSON(http.StatusInternalServerError, message)
				c.Abort()
				return
			}

			c.SetCookie(cookieName, token, 3600, "/", "", false, true)

			c.Status(http.StatusOK)
			c.Abort()
			return
		}

		userID, err := ssn.ValidateToken(cookie)
		if err != nil {
			sgr.Errorw(
				"failed to validate token",
				logKeyError, err.Error(),
			)

			message := models.Response{
				Message: http.StatusText(http.StatusUnauthorized),
			}

			c.JSON(http.StatusUnauthorized, message)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
	}
}

type gzipProvider struct {
	writer gin.ResponseWriter
	reader io.ReadCloser
}
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
func (p *gzipProvider) gzipHandler() *gzip.Writer {
	zw := gzip.NewWriter(p.writer)
	return zw
}
func (p *gzipProvider) unGzipHandler(sgr *zap.SugaredLogger) (*gzip.Reader, error) {
	zr, err := gzip.NewReader(p.reader)
	if err != nil {
		sgr.Errorw(
			"gzip middleware reader failed",
			logKeyError, err.Error(),
		)

		return nil, fmt.Errorf("failed to ungzip request body: %w", err)
	}

	return zr, nil
}
func GzipMiddleware(sgr *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.Request.Header.Get("Content-Type")
		supportsJSON := strings.Contains(contentType, "application/json")
		supportsHTML := strings.Contains(contentType, "text/html")

		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip && (supportsJSON || supportsHTML) {
			p := gzipProvider{
				writer: c.Writer,
			}
			zw := p.gzipHandler()
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
				ResponseWriter: p.writer,
				zw:             zw,
			}
			c.Writer.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			p := gzipProvider{
				reader: c.Request.Body,
			}
			zr, err := p.unGzipHandler(sgr)
			defer func() {
				err := zr.Close()
				if err != nil {
					sgr.Errorw(
						"gzip middleware reader close failed",
						logKeyError, err.Error(),
					)
				}
			}()

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
func LoggerMiddleware(sgr *zap.SugaredLogger) gin.HandlerFunc {
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

func TimeoutMiddleware(sgr *zap.SugaredLogger, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			sgr.Errorw(
				"failed to handle request by timeout",
				logKeyError, ctx.Err().Error(),
			)

			message := models.Response{
				Message: http.StatusText(http.StatusBadRequest),
			}

			c.JSON(http.StatusGatewayTimeout, message)
			c.Abort()
		}
	}
}
