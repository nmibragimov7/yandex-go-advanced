package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

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
			name:   "positive test #1",
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
			name:   "negative test #2",
			method: http.MethodGet,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:        405,
				contentType: "text/plain; charset=utf-8",
				response:    "Method Not Allowed",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, bytes.NewBuffer([]byte(test.body)))
			w := httptest.NewRecorder()
			MainPage(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.contentLength, res.Header.Get("Content-Length"))

			resBody, err := io.ReadAll(res.Body)
			defer res.Body.Close()

			require.NoError(t, err)
			assert.Contains(t, string(resBody), test.want.response)
		})
	}
}

func TestIdPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{
			name:   "positive test #1",
			method: http.MethodGet,
			want: want{
				code: 307,
			},
		},
		{
			name:   "negative test #2",
			method: http.MethodPost,
			want: want{
				code:        405,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://practicum.yandex.ru/")))
			w := httptest.NewRecorder()
			MainPage(w, request)

			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			defer res.Body.Close()

			parsedURL, err := url.Parse(string(resBody))
			require.NoError(t, err)

			request = httptest.NewRequest(test.method, parsedURL.Path, nil)
			w = httptest.NewRecorder()
			IDPage(w, request)

			res = w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
