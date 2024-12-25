package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipMiddleware(t *testing.T) {
	cnf := config.Init().GetConfig()
	sgr := logger.InitLogger()
	gzp := &GzipProvider{}
	lgp := &LoggerProvider{}
	str, err := storage.NewFileStorage(*cnf.FilePath)
	if err != nil {
		sgr.Errorw(
			"",
			"error", err.Error(),
		)
	}
	hdp := &handlers.HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}

	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()

	rtr := router.Provider{
		Config:           cnf,
		Storage:          str,
		Sugar:            sgr,
		GzipMiddleware:   gzp,
		LoggerMiddleWare: lgp,
		Handler:          hdp,
	}

	type wantShorten struct {
		code            int
		contentEncoding string
		contentType     string
	}
	testsShorten := []struct {
		name            string
		method          string
		path            string
		contentEncoding string
		headers         map[string]string
		body            models.ShortenRequestBody
		want            wantShorten
	}{
		{
			name:   "positive gzip middleware shorten api test #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			headers: map[string]string{
				"Accept-Encoding":  "gzip",
				"Content-Encoding": "gzip",
				"Content-Type":     "application/json",
			},
			body: models.ShortenRequestBody{URL: "https://practicum.yandex.ru/"},
			want: wantShorten{
				code:            http.StatusCreated,
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
		{
			name:   "negative gzip middleware shorten api test #2",
			method: http.MethodPost,
			path:   "/api/shorten",
			headers: map[string]string{
				"Accept-Encoding":  "gzip",
				"Content-Encoding": "gzip",
				"Content-Type":     "application/json",
			},
			body: models.ShortenRequestBody{URL: ""},
			want: wantShorten{
				code:            http.StatusBadRequest,
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
	}

	for _, test := range testsShorten {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(rtr.Router())
			defer ts.Close()

			bts, err := json.Marshal(test.body)
			assert.NoError(t, err)

			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err = zb.Write(bts)
			assert.NoError(t, err)

			err = zb.Close()
			assert.NoError(t, err)

			res, _ := util.TestRequest(t, ts, test.method, test.path, buf, test.headers)
			defer func() {
				err = res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentEncoding, res.Header.Get("Content-Encoding"))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

	type wantID struct {
		code int
	}
	testsID := []struct {
		name    string
		method  string
		headers map[string]string
		want    wantID
	}{
		{
			name:   "positive gzip middleware id api test #1",
			method: http.MethodGet,
			headers: map[string]string{
				"Accept-Encoding": "gzip",
			},
			want: wantID{
				code: http.StatusOK,
			},
		},
		{
			name:   "negative gzip middleware id api test #2",
			method: http.MethodPost,
			headers: map[string]string{
				"Accept-Encoding": "gzip",
			},
			want: wantID{
				code: http.StatusNotFound,
			},
		},
	}

	for _, test := range testsID {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(rtr.Router())
			defer ts.Close()

			headers := map[string]string{}
			resp, resBody := util.TestRequest(t, ts, http.MethodPost, "/", bytes.NewBufferString("https://google.kz/"), headers)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			parsedURL, err := url.Parse(string(resBody))
			require.NoError(t, err)
			require.NotNil(t, parsedURL, "Parsed URL should not be nil")

			res, _ := util.TestRequest(t, ts, test.method, parsedURL.Path, nil, test.headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
