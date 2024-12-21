package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	conf := config.Init()
	store := storage.NewStore()
	sugar := logger.InitLogger()
	mp := &Provider{}
	hp := &handlers.Provider{}

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
		body            models.ShortenRequestBody
		want            wantShorten
	}{
		{
			name:   "positive gzip middleware shorten api test #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   models.ShortenRequestBody{URL: "https://practicum.yandex.ru/"},
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
			body:   models.ShortenRequestBody{URL: ""},
			want: wantShorten{
				code:            http.StatusBadRequest,
				contentEncoding: "gzip",
				contentType:     "application/json",
			},
		},
	}

	for _, test := range testsShorten {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(router.Router(conf, store, sugar, mp, hp))
			defer ts.Close()

			bts, err := json.Marshal(test.body)
			assert.NoError(t, err)

			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err = zb.Write(bts)
			assert.NoError(t, err)

			err = zb.Close()
			assert.NoError(t, err)

			headers := map[string]string{
				"Accept-Encoding":  "gzip",
				"Content-Encoding": "gzip",
				"Content-Type":     "application/json",
			}
			res, _ := util.TestRequest(t, ts, test.method, test.path, buf, headers)
			defer func() {
				err = res.Body.Close()
				assert.NoError(t, err)
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentEncoding, res.Header.Get("Content-Encoding"))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
