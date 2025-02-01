package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/storage/db/shortener"
	"yandex-go-advanced/internal/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerProvider struct {
	Config  *config.Config
	Storage storage.Storage
	Sugar   *zap.SugaredLogger
	Session *session.SessionProvider
}

const (
	logKeyError     = "error"
	logKeyURI       = "uri"
	logKeyIP        = "ip"
	contentType     = "Content-Type"
	contentLength   = "Content-Length"
	applicationJSON = "application/json"
	shortenerTable  = "shortener"
)

func sendErrorResponse(c *gin.Context, sgr *zap.SugaredLogger, err error) {
	sgr.With(
		logKeyURI, c.Request.URL.Path,
		logKeyIP, c.ClientIP(),
	).Error(
		err,
	)

	message := models.Response{
		Message: http.StatusText(http.StatusInternalServerError),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		sgr.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)
		return
	}

	c.Header(contentType, applicationJSON)
	c.Header(contentLength, strconv.Itoa(len(bytes)))

	c.JSON(http.StatusInternalServerError, message)
}

func (p *HandlerProvider) MainPage(c *gin.Context) {
	var userID int64
	var err error
	if *p.Config.DataBase != "" {
		if userID, err = p.Session.ParseToken(c); err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)

			message := models.Response{
				Message: http.StatusText(http.StatusUnauthorized),
			}

			c.JSON(http.StatusUnauthorized, message)
			return
		}
	}

	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusMethodNotAllowed))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	url := string(body)

	key := util.GetKey()
	record := &models.ShortenRecord{
		ShortURL:    key,
		OriginalURL: url,
		UserID:      userID,
	}

	_, err = p.Storage.Set(shortenerTable, record)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)

		var duplicateError *shortener.DuplicateError
		if errors.As(err, &duplicateError) {
			c.Writer.WriteHeader(http.StatusConflict)
			c.Header(contentType, "text/plain")
			c.Header(contentLength, strconv.Itoa(len(*p.Config.BaseURL+"/"+duplicateError.ShortURL)))
			_, err = c.Writer.WriteString(*p.Config.BaseURL + "/" + duplicateError.ShortURL)
			if err != nil {
				p.Sugar.With(
					logKeyURI, c.Request.URL.Path,
					logKeyIP, c.ClientIP(),
				).Error(
					err,
				)
			}
			return
		}

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header(contentType, "text/plain")
	c.Header(contentLength, strconv.Itoa(len(*p.Config.BaseURL+"/"+key)))
	_, err = c.Writer.WriteString(*p.Config.BaseURL + "/" + key)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)
		return
	}
}
func (p *HandlerProvider) IDPage(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusMethodNotAllowed))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	path := c.Param("id")
	if path == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusNotFound))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	rec, err := p.Storage.Get(shortenerTable, path)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	record, ok := rec.(*models.ShortenRecord)
	if !ok {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			fmt.Errorf("record is not of type *models.ShortenRecord, actual type: %T", rec),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	if record.DeletedFlag {
		c.Writer.WriteHeader(http.StatusGone)
		return
	}

	if record.OriginalURL == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusNotFound))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
			return
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, record.OriginalURL)
}
func (p *HandlerProvider) ShortenHandler(c *gin.Context) {
	var userID int64
	var err error
	if *p.Config.DataBase != "" {
		if userID, err = p.Session.ParseToken(c); err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)

			message := models.Response{
				Message: http.StatusText(http.StatusUnauthorized),
			}

			c.JSON(http.StatusUnauthorized, message)
			return
		}
	}

	var body models.ShortenRequestBody
	bytes, err := c.GetRawData()
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}
	if err := json.Unmarshal(bytes, &body); err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	if body.URL == "" {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			fmt.Errorf("body url is empty: %s", body),
		)

		message := models.Response{
			Message: http.StatusText(http.StatusBadRequest),
		}

		c.Header(contentType, applicationJSON)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	key := util.GetKey()
	record := &models.ShortenRecord{
		ShortURL:    key,
		OriginalURL: body.URL,
		UserID:      userID,
	}

	_, err = p.Storage.Set(shortenerTable, record)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)

		var duplicateError *shortener.DuplicateError
		if errors.As(err, &duplicateError) {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Warn(
				err.Error(),
			)

			response := models.ShortenResponseSuccess{
				Result: *p.Config.BaseURL + "/" + duplicateError.ShortURL,
			}

			c.Header(contentType, applicationJSON)
			c.JSON(http.StatusConflict, response)
			return
		}

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)
		}
		return
	}

	response := models.ShortenResponseSuccess{
		Result: *p.Config.BaseURL + "/" + key,
	}

	c.Header(contentType, applicationJSON)
	c.JSON(http.StatusCreated, response)
}
func (p *HandlerProvider) PingHandler(c *gin.Context) {
	ctx := c.Request.Context()

	if p.Config.DataBase == nil {
		sendErrorResponse(c, p.Sugar, errors.New("database connection is nil"))
		return
	}

	err := p.Storage.Ping(ctx)
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "database is connected"})
}
func (p *HandlerProvider) ShortenBatchHandler(c *gin.Context) {
	var userID int64
	var err error
	if *p.Config.DataBase != "" {
		if userID, err = p.Session.ParseToken(c); err != nil {
			p.Sugar.With(
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			).Error(
				err,
			)

			message := models.Response{
				Message: http.StatusText(http.StatusUnauthorized),
			}

			c.JSON(http.StatusUnauthorized, message)
			return
		}
	}

	var body []models.ShortenBatchRequest
	bytes, err := c.GetRawData()
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}
	if err := json.Unmarshal(bytes, &body); err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	values := make([]interface{}, 0, len(body))
	result := make([]models.ShortenBatchResponse, 0, len(body))
	for _, value := range body {
		key := util.GetKey()
		values = append(values, &models.ShortenRecord{
			OriginalURL: value.OriginalURL,
			ShortURL:    key,
			UserID:      userID,
		})
		result = append(result, models.ShortenBatchResponse{
			CorrelationID: value.CorrelationID,
			ShortURL:      *p.Config.BaseURL + "/" + key,
		})
	}

	err = p.Storage.SetAll(shortenerTable, values)
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	c.Header(contentType, applicationJSON)
	c.JSON(http.StatusCreated, result)
}
func (p *HandlerProvider) UserUrlsHandler(c *gin.Context) {
	userID, err := p.Session.ParseToken(c)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)

		message := models.Response{
			Message: http.StatusText(http.StatusUnauthorized),
		}

		c.JSON(http.StatusUnauthorized, message)
		return
	}

	token, _ := c.Get("cookie")
	p.Sugar.Infow(
		"user ID info",
		"user_id", userID,
		"token", token,
	)

	rcs, err := p.Storage.GetAll(shortenerTable, userID)
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	if len(rcs) == 0 {
		message := models.Response{
			Message: http.StatusText(http.StatusNoContent),
		}

		c.JSON(http.StatusNoContent, message)
		return
	}

	records := make([]interface{}, 0, len(rcs))
	for _, rc := range rcs {
		value, ok := rc.(models.ShortenRecord)
		if !ok {
			sendErrorResponse(c, p.Sugar, errors.New("invalid shorten record"))
			return
		}
		records = append(records, map[string]interface{}{
			"short_url":    *p.Config.BaseURL + "/" + value.ShortURL,
			"original_url": value.OriginalURL,
		})
	}

	c.Header(contentType, applicationJSON)
	c.JSON(http.StatusOK, records)
}
func (p *HandlerProvider) UserUrlsDeleteHandler(c *gin.Context) {
	userID, err := p.Session.ParseToken(c)
	if err != nil {
		p.Sugar.With(
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		).Error(
			err,
		)

		message := models.Response{
			Message: http.StatusText(http.StatusUnauthorized),
		}

		c.JSON(http.StatusUnauthorized, message)
		return
	}

	token, _ := c.Get("cookie")
	p.Sugar.Infow(
		"user ID info",
		"user_id", userID,
		"token", token,
	)

	var body []string
	bytes, err := c.GetRawData()
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}
	if err := json.Unmarshal(bytes, &body); err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	generate := func(userID int64, key string) chan interface{} {
		out := make(chan interface{}, 1)
		go func() {
			defer close(out)
			val := &models.ShortenBatchUpdateRequest{
				ShortURL: key,
				UserID:   userID,
			}
			out <- val
		}()

		return out
	}

	values := make([]chan interface{}, 0, len(body))
	for _, value := range body {
		values = append(values, generate(userID, value))
	}

	go func() {
		done := make(chan struct{})
		defer close(done)
		p.Storage.AddToChannel(shortenerTable, done, values...)
	}()

	c.Status(http.StatusAccepted)
}
