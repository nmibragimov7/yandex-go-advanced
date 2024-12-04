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

func Router(store *storage.Store) *gin.Engine {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		MainPage(c, store)
	})
	r.GET("/:id", func(c *gin.Context) {
		IDPage(c, store)
	})

	return r
}

func MainPage(c *gin.Context, store *storage.Store) {
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

	store.SaveStore(key, url)

	conf := config.GetConfig()

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Length", fmt.Sprintf("%d", len(*conf.BaseURL+"/"+key)))
	_, err = c.Writer.Write([]byte(*conf.BaseURL + "/" + key))
	if err != nil {
		log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		return
	}
}

func IDPage(c *gin.Context, store *storage.Store) {
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

	if store.Store[path] == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.Write([]byte(http.StatusText(http.StatusNotFound)))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
			return
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, store.Store[path])
}
