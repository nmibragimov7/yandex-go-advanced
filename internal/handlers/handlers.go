package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/util"

	"github.com/gin-gonic/gin"
)

type Provider struct{}

const (
	logKeyError     = "error"
	logKeyURI       = "uri"
	logKeyIP        = "ip"
	contentType     = "Content-Type"
	contentLength   = "Content-Length"
	applicationJSON = "application/json"
)

func sendErrorResponse(c *gin.Context, sgr *logger.Logger, err error) {
	sugar := sgr.Get()

	sugar.Error(
		logKeyError, err.Error(),
		logKeyURI, c.Request.URL.Path,
		logKeyIP, c.ClientIP(),
	)

	message := models.ShortenResponseError{
		Message: http.StatusText(http.StatusInternalServerError),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		sugar.Error(
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

func (p *Provider) MainPage(c *gin.Context, cnf *config.Config, str *storage.FileStorage, sgr *logger.Logger) {
	sugar := sgr.Get()

	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusMethodNotAllowed))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	url := string(body)
	key := util.GetKey()

	result, err := str.ReadRecord()
	if err != nil {
		sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	uuid := 1
	if result != nil {
		uuid = result.UUID + 1
	}

	record := &models.ShortenRecord{
		UUID:        uuid,
		ShortURL:    key,
		OriginalURL: url,
	}
	err = str.WriteRecord(record)
	if err != nil {
		sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	configs := cnf.GetConfig()

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header(contentType, "text/plain")
	c.Header(contentLength, strconv.Itoa(len(*configs.BaseURL+"/"+key)))
	_, err = c.Writer.WriteString(*configs.BaseURL + "/" + key)
	if err != nil {
		sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)
		return
	}
}

func (p *Provider) IDPage(c *gin.Context, str *storage.FileStorage, sgr *logger.Logger) {
	sugar := sgr.Get()

	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusMethodNotAllowed))
		if err != nil {
			sugar.Error(
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
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	value := str.GetByKey(path)

	if value == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.WriteString(http.StatusText(http.StatusNotFound))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
			return
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, value)
}

func (p *Provider) ShortenHandler(c *gin.Context, cnf *config.Config, str *storage.FileStorage, sgr *logger.Logger) {
	sugar := sgr.Get()

	var body models.ShortenRequestBody
	bytes, err := c.GetRawData()
	if err != nil {
		sendErrorResponse(c, sgr, err)
		return
	}
	if err := json.Unmarshal(bytes, &body); err != nil {
		sendErrorResponse(c, sgr, err)
		return
	}

	if body.URL == "" {
		message := models.ShortenResponseError{
			Message: http.StatusText(http.StatusBadRequest),
		}

		c.Header(contentType, applicationJSON)
		c.JSON(http.StatusBadRequest, message)
		return
	}

	key := util.GetKey()

	result, err := str.ReadRecord()
	if err != nil {
		sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	uuid := 1
	if result != nil {
		uuid = result.UUID + 1
	}

	record := &models.ShortenRecord{
		UUID:        uuid,
		ShortURL:    key,
		OriginalURL: body.URL,
	}
	err = str.WriteRecord(record)
	if err != nil {
		sugar.Error(
			logKeyError, err.Error(),
			logKeyURI, c.Request.URL.Path,
			logKeyIP, c.ClientIP(),
		)

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.WriteString(http.StatusText(http.StatusInternalServerError))
		if err != nil {
			sugar.Error(
				logKeyError, err.Error(),
				logKeyURI, c.Request.URL.Path,
				logKeyIP, c.ClientIP(),
			)
		}
		return
	}

	configs := cnf.GetConfig()

	response := models.ShortenResponse{
		Result: *configs.BaseURL + "/" + key,
	}

	c.Header(contentType, applicationJSON)
	c.JSON(http.StatusCreated, response)
}
