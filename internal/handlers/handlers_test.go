package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/storage/handlersmocks"
	"yandex-go-advanced/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func createFormData(t *testing.T) (*bytes.Buffer, string, error) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("url", "https://practicum.yandex.ru/")
	if err != nil {
		return nil, "", err
	}

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func TestSendErrorResponse(t *testing.T) {
	sgr := zaptest.NewLogger(t).Sugar()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Real-IP", "127.0.0.1")

	sendErrorResponse(c, sgr, errors.New("test error"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, strconv.Itoa(len(w.Body.Bytes())), w.Header().Get("Content-Length"))

	var response models.Response
	err := json.NewDecoder(bytes.NewReader(w.Body.Bytes())).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), response.Message)
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
			name:   "test #1: status created",
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
			name:   "test #2: status conflict",
			method: http.MethodPost,
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				code:          http.StatusConflict,
				contentType:   "text/plain",
				contentLength: "30",
				response:      "http://localhost:8080/",
			},
		},
		{
			name:   "test #2: status not found",
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

	cnf := config.Init()
	sgr := logger.Init()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			"error", err.Error(),
		)
	}
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}
	ssp := &session.SessionProvider{
		Config: cnf,
	}
	rtp := router.RouterProvider{
		Storage: str,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	ts := httptest.NewServer(rtp.Router())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			headers := map[string]string{}
			res, resBody := util.TestRequest(t, ts, test.method, test.path, bytes.NewBufferString(test.body), headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("failed to close body: %s", err.Error())
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

	formData, contentType, err := createFormData(t)
	if err != nil {
		t.Fatalf("failed to create form-data: %s", err)
	}

	test := struct {
		name   string
		method string
		path   string
		body   *bytes.Buffer
		want   want
	}{
		name:   "test #1: status internal server error",
		method: http.MethodPost,
		path:   "/",
		body:   formData,
		want: want{
			code:          http.StatusInternalServerError,
			contentType:   "text/plain; charset=utf-8",
			contentLength: "21",
			response:      "Internal Server Error",
		},
	}

	t.Run(test.name, func(t *testing.T) {
		headers := map[string]string{
			"Content-Type": contentType,
		}
		res, resBody := util.TestRequest(t, ts, test.method, test.path, test.body, headers)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
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
			name:   "test #1: status ok",
			method: http.MethodGet,
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "test #2: status not found",
			method: http.MethodPost,
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	cnf := config.Init()
	sgr := logger.Init()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			"error", err.Error(),
		)
	}
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}
	ssp := &session.SessionProvider{
		Config: cnf,
	}
	rtp := router.RouterProvider{
		Storage: str,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	ts := httptest.NewServer(rtp.Router())
	defer ts.Close()

	headers := map[string]string{}
	resp, resBody := util.TestRequest(t, ts, http.MethodPost, "/", bytes.NewBufferString("https://google.kz/"), headers)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close body: %s", err.Error())
		}
	}()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsedURL, err := url.Parse(string(resBody))
			require.NoError(t, err)
			require.NotNil(t, parsedURL, "Parsed URL should not be nil")

			headers = map[string]string{}
			res, _ := util.TestRequest(t, ts, test.method, parsedURL.Path, nil, headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("failed to close body: %s", err.Error())
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
			name:   "test #1: status created",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   models.ShortenRequestBody{URL: "https://practicum.yandex.kz/"},
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:   "test #2: status conflict",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   models.ShortenRequestBody{URL: "https://practicum.yandex.kz/"},
			want: want{
				code:        http.StatusConflict,
				contentType: "application/json",
			},
		},
		{
			name:   "test #3: status bad request",
			method: http.MethodPost,
			path:   "/api/shorten",
			body:   models.ShortenRequestBody{URL: ""},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}

	cnf := config.Init()
	sgr := logger.Init()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			"error", err.Error(),
		)
	}
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}
	ssp := &session.SessionProvider{
		Config: cnf,
	}
	rtp := router.RouterProvider{
		Storage: str,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(rtp.Router())
			defer ts.Close()

			bts, err := json.Marshal(test.body)
			assert.NoError(t, err)

			buf := bytes.NewBuffer(bts)
			headers := map[string]string{}
			res, _ := util.TestRequest(t, ts, test.method, test.path, buf, headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("failed to close body: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestPingHandler(t *testing.T) {
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
			name:   "test #1: status ok",
			method: http.MethodGet,
			path:   "/ping",
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}

	cnf := config.Init()
	sgr := logger.Init()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			"error", err.Error(),
		)
	}
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}
	ssp := &session.SessionProvider{
		Config: cnf,
	}
	rtp := router.RouterProvider{
		Storage: str,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(rtp.Router())
			defer ts.Close()

			headers := map[string]string{}
			res, _ := util.TestRequest(t, ts, test.method, test.path, nil, headers)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("failed to close body: %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func ptr(t *testing.T, s string) *string {
	t.Helper()

	return &s
}

func getMockInitial(t *testing.T) (*gomock.Controller, *gin.Engine, *handlersmocks.MockStorage) {
	ctrl := gomock.NewController(t)
	mockStorage := handlersmocks.NewMockStorage(ctrl)

	cnf := &config.Config{
		Server:   ptr(t, ":8080"),
		BaseURL:  ptr(t, "http://localhost:8080"),
		FilePath: ptr(t, "./storage.txt"),
		DataBase: ptr(t, "host=localhost user=postgres password=admin dbname=postgres sslmode=disable"),
	}
	sgr := logger.Init()
	ssp := &session.SessionProvider{
		Config: cnf,
	}
	hdp := &HandlerProvider{
		Config:  cnf,
		Storage: mockStorage,
		Sugar:   sgr,
		Session: ssp,
	}
	rtp := router.RouterProvider{
		Storage: mockStorage,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	return ctrl, rtp.Router(), mockStorage
}

func TestShortenBatchHandler(t *testing.T) {
	ctrl, engine, mockStorage := getMockInitial(t)
	defer ctrl.Finish()

	t.Run("Save shortener batch", func(t *testing.T) {
		body := `[{"correlation_id":"123","original_url":"http://practicum.yandex.ru"}]`

		mockStorage.EXPECT().SetAll("shortener", strings.NewReader(body)).AnyTimes()

		ssp := session.SessionProvider{
			Config: nil,
		}
		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)

		cookie := &http.Cookie{
			Name:  "user_token",
			Value: token,
			Path:  "/",
		}

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL+"/api/user/urls", strings.NewReader(`["123456","234567","345678"]`))
		require.NoError(t, err)

		req.AddCookie(cookie)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusAccepted, res.StatusCode)
	})

	t.Run("Save shortener batch with status internal server error", func(t *testing.T) {
		mockStorage.EXPECT().AddToChannel("shortener", gomock.Any(), gomock.Any()).AnyTimes()

		ssp := session.SessionProvider{
			Config: nil,
		}
		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)

		cookie := &http.Cookie{
			Name:  "user_token",
			Value: token,
			Path:  "/",
		}

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL+"/api/user/urls", nil)
		require.NoError(t, err)

		req.AddCookie(cookie)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestUserUrlsHandler(t *testing.T) {
	ctrl, engine, mockStorage := getMockInitial(t)
	defer ctrl.Finish()

	shortenRecords := []models.ShortenRecord{
		{ShortURL: "123456", OriginalURL: "https://practicum.yandex.kz/"},
	}
	recordsAsInterface := make([]interface{}, len(shortenRecords))
	for i, v := range shortenRecords {
		recordsAsInterface[i] = v
	}

	t.Run("Get user urls", func(t *testing.T) {
		mockStorage.EXPECT().Set("users", &models.UserRecord{}).Return(int64(1), nil)
		mockStorage.EXPECT().GetAll("shortener", int64(1)).Return(recordsAsInterface, nil)

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL+"/api/user/urls", nil)
		require.NoError(t, err)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.JSONEq(t, `[{"original_url":"https://practicum.yandex.kz/","short_url":"http://localhost:8080/123456"}]`, string(resBody))
	})

	t.Run("Get user urls with status unauthorized", func(t *testing.T) {
		mockStorage.EXPECT().Set("users", &models.UserRecord{}).Return(int64(0), nil)

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL+"/api/user/urls", nil)
		require.NoError(t, err)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Get user urls with no content", func(t *testing.T) {
		mockStorage.EXPECT().Set("users", &models.UserRecord{}).Return(int64(1), nil)
		mockStorage.EXPECT().GetAll("shortener", int64(1)).Return([]interface{}{}, nil)

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL+"/api/user/urls", nil)
		require.NoError(t, err)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}

func TestUserUrlsDeleteHandler(t *testing.T) {
	ctrl, engine, mockStorage := getMockInitial(t)
	defer ctrl.Finish()

	t.Run("Delete my urls", func(t *testing.T) {
		mockStorage.EXPECT().AddToChannel("shortener", gomock.Any(), gomock.Any()).AnyTimes()

		ssp := session.SessionProvider{
			Config: nil,
		}
		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)

		cookie := &http.Cookie{
			Name:  "user_token",
			Value: token,
			Path:  "/",
		}

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL+"/api/user/urls", strings.NewReader(`["123456","234567","345678"]`))
		require.NoError(t, err)

		req.AddCookie(cookie)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusAccepted, res.StatusCode)
	})

	t.Run("Delete my urls with status unauthorized", func(t *testing.T) {
		mockStorage.EXPECT().AddToChannel("shortener", gomock.Any(), gomock.Any()).AnyTimes()

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL+"/api/user/urls", strings.NewReader(`["123456","234567","345678"]`))
		require.NoError(t, err)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("Delete my urls with status internal server error", func(t *testing.T) {
		mockStorage.EXPECT().AddToChannel("shortener", gomock.Any(), gomock.Any()).AnyTimes()

		ssp := session.SessionProvider{
			Config: nil,
		}
		token, err := ssp.GenerateToken(int64(1))
		assert.NoError(t, err)

		cookie := &http.Cookie{
			Name:  "user_token",
			Value: token,
			Path:  "/",
		}

		ts := httptest.NewServer(engine)
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL+"/api/user/urls", nil)
		require.NoError(t, err)

		req.AddCookie(cookie)

		res, err := ts.Client().Do(req)
		require.NoError(t, err)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Printf("failed to close body: %s", err.Error())
			}
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
