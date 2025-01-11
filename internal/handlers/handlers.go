package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerProvider struct {
	Config  *config.Config
	Storage storage.Storage
	Sugar   *zap.SugaredLogger
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
	sgr.Error(
		logKeyError, err.Error(),
		logKeyURI, c.Request.URL.Path,
		logKeyIP, c.ClientIP(),
	)

	message := models.Response{
		Message: http.StatusText(http.StatusInternalServerError),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		sgr.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)
		return
	}

	c.Header(contentType, applicationJSON)
	c.Header(contentLength, strconv.Itoa(len(bytes)))

	c.JSON(http.StatusInternalServerError, message)
}

func (p *HandlerProvider) MainPage(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusMethodNotAllowed))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		p.Sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	url := string(body)
	key := util.GetKey()
	record := &models.ShortenRecord{
		ShortURL:    key,
		OriginalURL: url,
	}

	err = p.Storage.Set(shortenerTable, record)
	if err != nil {
		p.Sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header(contentType, "text/plain")
	c.Header(contentLength, strconv.Itoa(len(*p.Config.BaseURL+"/"+key)))
	_, err = c.Writer.WriteString(*p.Config.BaseURL + "/" + key)
	if err != nil {
		p.Sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)
		return
	}
}
func (p *HandlerProvider) IDPage(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusMethodNotAllowed))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	path := c.Param("id")
	if path == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusNotFound))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	rec, err := p.Storage.Get(shortenerTable, path)
	if err != nil {
		p.Sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	record, ok := rec.(*models.ShortenRecord)
	if !ok {
		p.Sugar.Error(
			logKeyError, errors.New("record is not of type *models.ShortenRecord"),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	if record.OriginalURL == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusNotFound))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
			return
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, record.OriginalURL)
}
func (p *HandlerProvider) ShortenHandler(c *gin.Context) {
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
	}

	err = p.Storage.Set(shortenerTable, record)
	if err != nil {
		p.Sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

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
		})
		result = append(result, models.ShortenBatchResponse{
			CorrelationID: value.CorrelationID,
			ShortURL:      *p.Config.BaseURL + "/" + key,
		})
	}

	err = p.Storage.SetByTransaction(shortenerTable, values)
	if err != nil {
		p.Sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			p.Sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	c.Header(contentType, applicationJSON)
	c.JSON(http.StatusCreated, result)
}
