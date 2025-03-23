package util

import (
	"bytes"
	"crypto/rand"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	var shortID strings.Builder

	shortID.Grow(length)

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes")
	}

	for _, b := range randomBytes {
		shortID.WriteByte(charset[b%byte(len(charset))])
	}

	return shortID.String()
}

func TestRequest(
	t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	body *bytes.Buffer,
	headers map[string]string,
) (*http.Response, []byte) {
	t.Helper()

	if body == nil {
		body = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("failed to close body: %s", err.Error())
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return res, resBody
}
