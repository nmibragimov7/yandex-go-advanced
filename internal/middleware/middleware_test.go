package middleware

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"
	"yandex-go-advanced/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type readCloser struct {
	reader io.Reader
}

func (rc *readCloser) Read(p []byte) (int, error) {
	return rc.reader.Read(p)
}

func (rc *readCloser) Close() error {
	// Ничего не делает, так как bytes.Reader не требует явного закрытия
	return nil
}

func TestGzipMiddleware(t *testing.T) {
	sgr := logger.InitLogger()

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	provider := gzipProvider{
		writer: context.Writer,
	}

	zw := provider.gzipHandler()
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

	provider = gzipProvider{
		reader: &readCloser{
			reader: bytes.NewReader(result),
		},
	}
	zr, err := provider.unGzipHandler(sgr)
	assert.NoError(t, err)
	defer func() {
		err := zr.Close()
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
}
