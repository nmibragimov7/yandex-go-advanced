package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/db"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name:   "positive main page test #1",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:          http.StatusCreated,
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
				code:          http.StatusNotFound,
				contentType:   "text/plain",
				contentLength: "18",
				response:      "404 page not found",
			},
		},
	}

	cnf := config.Init().GetConfig()
	sgr := logger.InitLogger()
	str, err := storage.InitFileStorage(*cnf.FilePath)
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()
	if err != nil {
		sgr.Errorw(
			"",
			"error", err.Error(),
		)
	}

	dbp := db.DatabaseProvider{
		Sugar:  sgr,
		Config: cnf,
	}
	database, err := dbp.Init()
	if err != nil {
		sgr.Errorw(
			"failed to init database",
			logKeyError, err.Error(),
		)
	}
	defer func() {
		if database != nil {
			err := database.Close()
			if err != nil {
				sgr.Errorw(
					"Failed to close database connection",
					logKeyError, err.Error(),
				)
			}
		}
	}()

	hdp := &HandlerProvider{
		Config:   cnf,
		Storage:  str,
		Sugar:    sgr,
		Database: database,
	}
	rtr := router.RouterProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		Handler: hdp,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(rtr.Router())
			defer ts.Close()

			headers := map[string]string{}
			res, resBody := util.TestRequest(t, ts, test.method, test.path, bytes.NewBufferString(test.body), headers)
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
				code: http.StatusOK,
			},
		},
		{
			name:   "negative id page test #2",
			method: http.MethodPost,
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	cnf := config.Init().GetConfig()
	sgr := logger.InitLogger()
	str, err := storage.InitFileStorage(*cnf.FilePath)
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()
	if err != nil {
		sgr.Errorw(
			"",
			"error", err.Error(),
		)
	}

	dbp := db.DatabaseProvider{
		Sugar:  sgr,
		Config: cnf,
	}
	database, err := dbp.Init()
	if err != nil {
		sgr.Errorw(
			"failed to init database",
			logKeyError, err.Error(),
		)
	}
	defer func() {
		if database != nil {
			err := database.Close()
			if err != nil {
				sgr.Errorw(
					"Failed to close database connection",
					logKeyError, err.Error(),
				)
			}
		}
	}()

	hdp := &HandlerProvider{
		Config:   cnf,
		Storage:  str,
		Sugar:    sgr,
		Database: database,
	}
	rtr := router.RouterProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		Handler: hdp,
	}

	for _, test := range tests {
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

			headers = map[string]string{}
			res, _ := util.TestRequest(t, ts, test.method, parsedURL.Path, nil, headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
func TestShortenHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   models.ShortenRequestBody
		want   want
	}{
		{
			name:   "positive shorten handler test #1",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   models.ShortenRequestBody{URL: "https://practicum.yandex.ru/"},
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:   "negative shorten handler test #2",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   models.ShortenRequestBody{URL: ""},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}

	cnf := config.Init().GetConfig()
	sgr := logger.InitLogger()
	str, err := storage.InitFileStorage(*cnf.FilePath)
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()
	if err != nil {
		sgr.Errorw(
			"",
			"error", err.Error(),
		)
	}

	dbp := db.DatabaseProvider{
		Sugar:  sgr,
		Config: cnf,
	}
	database, err := dbp.Init()
	if err != nil {
		sgr.Errorw(
			"failed to init database",
			logKeyError, err.Error(),
		)
	}
	defer func() {
		if database != nil {
			err := database.Close()
			if err != nil {
				sgr.Errorw(
					"Failed to close database connection",
					logKeyError, err.Error(),
				)
			}
		}
	}()

	hdp := &HandlerProvider{
		Config:   cnf,
		Storage:  str,
		Sugar:    sgr,
		Database: database,
	}
	rtr := router.RouterProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		Handler: hdp,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(rtr.Router())
			defer ts.Close()

			bts, err := json.Marshal(test.body)
			assert.NoError(t, err)

			buf := bytes.NewBuffer(bts)
			headers := map[string]string{}
			res, _ := util.TestRequest(t, ts, test.method, test.path, buf, headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("Response body close: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
