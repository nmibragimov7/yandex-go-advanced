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
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type HandlerProvider struct {
	Config   *config.Config
	Storage  storage.Storage
	Sugar    *zap.SugaredLogger
	Database *sqlx.DB
}

const (
	logKeyError     = "error"
	logKeyURI       = "uri"
	logKeyIP        = "ip"
	contentType     = "Content-Type"
	contentLength   = "Content-Length"
	applicationJSON = "application/json"
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

	if *p.Config.DataBase != "" {
		query := "INSERT INTO shortener (short_url, original_url) VALUES ($1, $2)"
		_, err = p.Database.Exec(query, record.ShortURL, record.OriginalURL)
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
	} else {
		err = p.Storage.Set(record)
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

	var value string
	var err error
	if *p.Config.DataBase != "" {
		var record models.ShortenRecord
		query := "SELECT short_url, original_url FROM shortener WHERE short_url = $1"
		err = p.Database.QueryRow(query, path).Scan(&record.ShortURL, &record.OriginalURL)
		value = record.OriginalURL
	} else {
		value, err = p.Storage.Get(path)
	}
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

	if value == "" {
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

	c.Redirect(http.StatusTemporaryRedirect, value)
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

	if *p.Config.DataBase != "" {
		query := "INSERT INTO shortener (short_url, original_url) VALUES ($1, $2)"
		_, err = p.Database.Exec(query, record.ShortURL, record.OriginalURL)
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
	} else {
		err = p.Storage.Set(record)
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

	if p.Database == nil {
		sendErrorResponse(c, p.Sugar, errors.New("database connection is nil"))
		return
	}

	err := p.Database.PingContext(ctx)
	if err != nil {
		sendErrorResponse(c, p.Sugar, err)
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "database is connected"})
}
