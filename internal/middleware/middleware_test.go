package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"yandex-go-advanced/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGzipHandler(t *testing.T) {
	sgr := logger.Init()

	t.Run("test gzip handler", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(recorder)

		gzp := gzipProvider{
			writer: context.Writer,
		}

		zw := gzp.gzipHandler()
		defer func() {
			err := zw.Close()
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()

		body := []byte("http://localhost:8080/735P2s38")
		_, err := zw.Write(body)
		assert.NoError(t, err)

		err = zw.Close()
		assert.NoError(t, err)

		result := recorder.Body.Bytes()
		assert.NotNil(t, result)
		assert.NotEqual(t, body, result)

		zr, err := gzip.NewReader(bytes.NewReader(result))
		assert.NoError(t, err)
		defer func() {
			err := zr.Close()
			assert.NoError(t, err)
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()

		decompressed, err := io.ReadAll(zr)
		assert.NoError(t, err)
		assert.Equal(t, body, decompressed)
	})
}

func TestUnGzipHandler(t *testing.T) {
	sgr := logger.Init()

	t.Run("test ungzip handler", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(recorder)

		uzp := gzipProvider{
			writer: context.Writer,
		}

		zw := uzp.gzipHandler()
		defer func() {
			err := zw.Close()
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()

		body := []byte("http://localhost:8080/735P2s38")
		_, err := zw.Write(body)
		assert.NoError(t, err)

		err = zw.Close()
		assert.NoError(t, err)

		result := recorder.Body.Bytes()
		assert.NotNil(t, result)
		assert.NotEqual(t, body, result)

		uzp.reader = io.NopCloser(bytes.NewReader(result))

		zr, err := uzp.unGzipHandler(sgr)
		assert.NoError(t, err)
		defer func() {
			err := zr.Close()
			assert.NoError(t, err)
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()

		decompressed, err := io.ReadAll(zr)
		assert.NoError(t, err)
		assert.Equal(t, body, decompressed)
	})
}

func TestGzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sgr := logger.Init()
	router := gin.New()
	router.Use(GzipMiddleware(sgr))

	router.POST("/", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read body"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": string(body)})
	})

	t.Run("request and response without gzip", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"message":"test"}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Header().Get("Content-Encoding"), "")
		assert.JSONEq(t, `{"received":"{\"message\":\"test\"}"}`, recorder.Body.String())
	})

	t.Run("response with gzip", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"message":"test"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip")

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "gzip", recorder.Header().Get("Content-Encoding"))

		zr, err := gzip.NewReader(recorder.Body)
		assert.NoError(t, err)
		defer func() {
			err := zr.Close()
			assert.NoError(t, err)
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()

		decompressed, err := io.ReadAll(zr)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"received":"{\"message\":\"test\"}"}`, string(decompressed))
	})
}

func TestLoggerMiddleware(t *testing.T) {
	sgr := logger.Init()

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://chatgpt.com"))

	mw := LoggerMiddleware(sgr)
	mw(context)

	assert.Equal(t, http.StatusOK, recorder.Code)

	size := recorder.Body.Len()
	assert.GreaterOrEqual(t, size, 0)
}

func TestTimeoutMiddleware(t *testing.T) {
	sgr := logger.Init()

	t.Run("timeout middleware success test", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(recorder)
		context.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://gitlab.com"))

		router := gin.New()
		router.Use(TimeoutMiddleware(sgr, 2*time.Second))
		router.POST("/", func(c *gin.Context) {
			select {
			case <-c.Request.Context().Done():
				return
			case <-time.After(1 * time.Second):
			}
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		router.ServeHTTP(recorder, context.Request)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("timeout middleware error test", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(recorder)
		context.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://github.com"))

		router := gin.New()
		router.Use(TimeoutMiddleware(sgr, 1*time.Millisecond))
		router.POST("/", func(c *gin.Context) {
			select {
			case <-c.Request.Context().Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		router.ServeHTTP(recorder, context.Request)

		assert.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	})
}
