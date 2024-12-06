package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/util"
)

func Router(cnf *config.Config, str *storage.Store) *gin.Engine {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		MainPage(c, cnf, str)
	})
	r.GET("/:id", func(c *gin.Context) {
		IDPage(c, str)
	})

	return r
}

func MainPage(c *gin.Context, cnf *config.Config, str *storage.Store) {
	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("ERROR: failed to read request body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	url := string(body)
	key := util.GetKey()

	str.SaveStore(key, url)

	configs := cnf.GetConfig()

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Length", fmt.Sprintf("%d", len(*configs.BaseURL+"/"+key)))
	_, err = c.Writer.Write([]byte(*configs.BaseURL + "/" + key))
	if err != nil {
		log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		return
	}
}

func IDPage(c *gin.Context, str *storage.Store) {
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	path := c.Param("id")
	if path == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.Write([]byte(http.StatusText(http.StatusNotFound)))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	if str.Store[path] == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.Write([]byte(http.StatusText(http.StatusNotFound)))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
			return
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, str.Store[path])
}
