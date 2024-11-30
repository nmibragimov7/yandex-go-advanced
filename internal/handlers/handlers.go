package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sync"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/util"
)

type Store struct {
	store map[string]string
	mtx   *sync.Mutex
}

func (s *Store) GetStore() map[string]string {
	return s.store
}
func (s *Store) SaveStore(key, url string) {
	globalStore.mtx.Lock()
	globalStore.store[key] = url
	globalStore.mtx.Unlock()
}

var globalStore = Store{
	store: make(map[string]string),
	mtx:   &sync.Mutex{},
}

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/", MainPage)
	r.GET("/:id", IDPage)

	return r
}

func MainPage(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.Write([]byte("Method Not Allowed"))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("ERROR: failed to read request body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())

		c.Writer.WriteHeader(http.StatusInternalServerError)
		_, err = c.Writer.Write([]byte("Failed to read request body"))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	url := string(body)
	key := util.GetKey()

	globalStore.SaveStore(key, url)

	baseURL := config.GetBaseURL()

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Length", fmt.Sprintf("%d", len(*baseURL+"/"+key)))
	_, err = c.Writer.Write([]byte(*baseURL + "/" + key))
	if err != nil {
		log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
	}
}

func IDPage(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := c.Writer.Write([]byte("Method Not Allowed"))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	path := c.Param("id")
	if path == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.Write([]byte("Not Found"))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	if globalStore.store[path] == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		_, err := c.Writer.Write([]byte("Not Found"))
		if err != nil {
			log.Printf("ERROR: failed to send response body: %v, Path: %s, IP: %s", err, c.Request.URL.Path, c.ClientIP())
		}
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, globalStore.store[path])
}
