package util

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetKey() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)[:8]
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
			log.Printf("Response body close: %s", err.Error())
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return res, resBody
}
