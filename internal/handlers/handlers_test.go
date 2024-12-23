package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	body *bytes.Buffer,
) (*http.Response, []byte) {
	t.Helper()

	if body == nil {
		body = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

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

func TestMainPage(t *testing.T) {
	type want struct {
		code          int
		contentType   string
		contentLength string
		response      string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			name:   "positive main page test #1",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:          201,
				contentType:   "text/plain",
				contentLength: "30",
				response:      "http://localhost:8080/",
			},
		},
		{
			name:   "negative main page test #2",
			method: http.MethodGet,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:          404,
				contentType:   "text/plain",
				contentLength: "18",
				response:      "404 page not found",
			},
		},
	}

	conf := config.Init()
	store := storage.NewStore()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(Router(conf, store))
			defer ts.Close()

			res, resBody := testRequest(t, ts, test.method, test.path, bytes.NewBufferString(test.body))
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			if test.want.contentLength != "" {
				assert.Equal(t, test.want.contentLength, res.Header.Get("Content-Length"))
			}
			assert.Contains(t, string(resBody), test.want.response)
		})
	}
}

func TestIdPage(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "positive id page test #1",
			method: http.MethodGet,
			want: want{
				code: 200,
			},
		},
		{
			name:   "negative id page test #2",
			method: http.MethodPost,
			want: want{
				code: 404,
			},
		},
	}

	conf := config.Init()
	store := storage.NewStore()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(Router(conf, store))
			defer ts.Close()

			resp, resBody := testRequest(t, ts, http.MethodPost, "/", bytes.NewBufferString("https://google.kz/"))
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			parsedURL, err := url.Parse(string(resBody))
			require.NoError(t, err)
			require.NotNil(t, parsedURL, "Parsed URL should not be nil")

			res, _ := testRequest(t, ts, test.method, parsedURL.Path, nil)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
