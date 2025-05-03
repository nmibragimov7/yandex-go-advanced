package util

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkGetKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetKey()
	}
}

func TestGetKey(t *testing.T) {
	t.Run("GetKey test", func(t *testing.T) {
		key := GetKey()
		assert.Equal(t, len(key), 8)
	})

	t.Run("GetKey test with error", func(t *testing.T) {
		orig := readRandomBytes
		defer func() { readRandomBytes = orig }()

		readRandomBytes = func(_ []byte) (int, error) {
			return 0, errors.New("mock error")
		}

		key := GetKey()
		assert.Equal(t, "", key)
	})
}

func TestTestRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			return
		}
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, body := TestRequest(t, ts, "GET", "/", nil, nil)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("failed to close body: %s", err.Error())
		}
	}()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "ok", string(body))
}
