package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockStorageSuccess struct{}
type MockStorageError struct{}

func (p *MockStorageSuccess) Get(_ string, _ string) (interface{}, error)                   { return nil, nil }
func (p *MockStorageSuccess) GetAll(_ string, _ interface{}) ([]interface{}, error)         { return nil, nil }
func (p *MockStorageSuccess) Set(_ string, _ interface{}) (interface{}, error)              { return int64(1), nil }
func (p *MockStorageSuccess) SetAll(_ string, _ []interface{}) error                        { return nil }
func (p *MockStorageSuccess) AddToChannel(_ string, _ chan struct{}, _ ...chan interface{}) {}
func (p *MockStorageSuccess) Close() error                                                  { return nil }
func (p *MockStorageSuccess) Ping(_ context.Context) error                                  { return nil }

func (p *MockStorageError) Get(_ string, _ string) (interface{}, error)           { return nil, nil }
func (p *MockStorageError) GetAll(_ string, _ interface{}) ([]interface{}, error) { return nil, nil }
func (p *MockStorageError) Set(_ string, _ interface{}) (interface{}, error) {
	return nil, errors.New("set error")
}
func (p *MockStorageError) SetAll(_ string, _ []interface{}) error                        { return nil }
func (p *MockStorageError) AddToChannel(_ string, _ chan struct{}, _ ...chan interface{}) {}
func (p *MockStorageError) Close() error                                                  { return nil }
func (p *MockStorageError) Ping(_ context.Context) error                                  { return nil }

type MockSession struct {
	*config.Config
}

func (p *MockSession) GenerateToken(_ int64) (string, error) {
	return "", errors.New("generate error")
}
func (p *MockSession) ParseCookie(_ string) (int64, error) {
	return int64(0), errors.New("check error")
}
func (p *MockSession) CheckCookie(_ string) error { return errors.New("check error") }

func TestAuthMiddleware(t *testing.T) {

	t.Run("auth middleware success test", func(t *testing.T) {
		cnf := config.Init()
		sgr := logger.Init()
		ssp := &session.SessionProvider{
			Config: cnf,
		}

		atp := &AuthProvider{
			Storage: nil,
			Config:  cnf,
			Sugar:   sgr,
			Session: ssp,
		}

		fnc := AuthMiddleware(atp)

		assert.NotNil(t, fnc)
	})

	err := os.Setenv("DATABASE_DSN", "postgres://localhost:5432/testdb")
	assert.NoError(t, err)

	cnf := config.Init()
	sgr := logger.Init()
	str := &MockStorageSuccess{}
	ssp := &MockSession{
		Config: cnf,
	}

	t.Run("auth middleware generate token error test", func(t *testing.T) {
		atp := &AuthProvider{
			Storage: str,
			Config:  cnf,
			Sugar:   sgr,
			Session: ssp,
		}

		r := gin.Default()
		r.Use(AuthMiddleware(atp))
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("auth middleware check token error test", func(t *testing.T) {
		atp := &AuthProvider{
			Storage: str,
			Config:  cnf,
			Sugar:   sgr,
			Session: ssp,
		}

		r := gin.Default()
		r.Use(AuthMiddleware(atp))
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		cookieValue := "your-cookie-value"
		req.AddCookie(&http.Cookie{
			Name:  cookieName,
			Value: cookieValue,
			Path:  "/",
		})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("auth middleware set user error test", func(t *testing.T) {
		strerr := &MockStorageError{}
		atp := &AuthProvider{
			Storage: strerr,
			Config:  cnf,
			Sugar:   sgr,
			Session: ssp,
		}

		r := gin.Default()
		r.Use(AuthMiddleware(atp))
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGzipHandler(t *testing.T) {
	sgr := logger.Init()

	t.Run("test gzip handler", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		gzp := gzipProvider{
			writer: ctx.Writer,
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
			err = zr.Close()
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
		ctx, _ := gin.CreateTestContext(recorder)

		uzp := gzipProvider{
			writer: ctx.Writer,
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
			err = zr.Close()
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

	t.Run("test ungzip handler with error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)

		uzp := gzipProvider{
			writer: ctx.Writer,
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

		gzipNewReader = func(r io.Reader) (*gzip.Reader, error) {
			return nil, errors.New("reader error")
		}
		_, err = uzp.unGzipHandler(sgr)
		assert.ErrorContains(t, err, "reader error")

		gzipNewReader = gzip.NewReader
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
			err = zr.Close()
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

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	t.Run("request with gzip", func(t *testing.T) {
		gzp := gzipProvider{
			writer: ctx.Writer,
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

		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(result))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		zr, err := gzip.NewReader(recorder.Body)
		assert.NoError(t, err)
		defer func() {
			err = zr.Close()
			assert.NoError(t, err)
			if err != nil {
				sgr.Errorw(
					"",
					"error", err.Error(),
				)
			}
		}()

		decompressed, err := io.ReadAll(zr)
		assert.Equal(t, string(body), string(decompressed))
	})

	t.Run("request with gzip with error", func(t *testing.T) {
		gzipNewReader = func(r io.Reader) (*gzip.Reader, error) {
			return nil, errors.New("reader error")
		}

		gzp := gzipProvider{
			writer: ctx.Writer,
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

		req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(result))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")

		router.ServeHTTP(recorder, req)

		assert.Contains(t, recorder.Body.String(), "Bad Request")
		gzipNewReader = gzip.NewReader
	})
}

func TestLoggerMiddleware(t *testing.T) {
	sgr := logger.Init()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://chatgpt.com"))

	mw := LoggerMiddleware(sgr)
	mw(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)

	size := recorder.Body.Len()
	assert.GreaterOrEqual(t, size, 0)
}

func TestTimeoutMiddleware(t *testing.T) {
	sgr := logger.Init()

	t.Run("timeout middleware success test", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://gitlab.com"))

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

		router.ServeHTTP(recorder, ctx.Request)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("timeout middleware error test", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://github.com"))

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

		router.ServeHTTP(recorder, ctx.Request)

		assert.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	})
}

type mockResponseWriter struct {
	gin.ResponseWriter
	statusCode int
	body       []byte
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.body = append(m.body, data...)
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestWrite(t *testing.T) {
	t.Run("request with gzip with error", func(t *testing.T) {
		expectedData := []byte("test data")

		mockWriter := &mockResponseWriter{}
		lgr := &loggerWriter{
			ResponseWriter: mockWriter,
			body:           bytes.NewBuffer([]byte{}),
		}

		n, err := lgr.Write(expectedData)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if n != len(expectedData) {
			t.Errorf("expected %d bytes written, got %d", len(expectedData), n)
		}

		if lgr.body.String() != string(expectedData) {
			t.Errorf("expected body to contain %s, got %s", expectedData, lgr.body.Bytes())
		}

		if string(mockWriter.body) != string(expectedData) {
			t.Errorf("expected mockWriter body to contain %s, got %s", expectedData, mockWriter.body)
		}
	})

}
