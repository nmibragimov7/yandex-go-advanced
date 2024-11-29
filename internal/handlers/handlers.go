package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/util"
)

var store = make(map[string]string)
var mtx sync.Mutex

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/", MainPage)
	r.GET("/:id", IDPage)

	//r := chi.NewRouter()
	//r.Use(CustomMiddleware)
	//
	//r.Post("/", MainPage)
	//r.Get(`/{id}`, IDPage)

	return r
}

func MainPage(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		c.Writer.Write([]byte("Method Not Allowed"))
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte("Failed to read request body"))
		return
	}

	url := string(body)
	key := util.GetKey()

	mtx.Lock()
	store[key] = url
	mtx.Unlock()

	c.Writer.WriteHeader(http.StatusCreated)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Length", fmt.Sprintf("%d", len(*config.BaseURL+"/"+key)))
	c.Writer.Write([]byte(*config.BaseURL + "/" + key))
}

func IDPage(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		c.Writer.Write([]byte("Method Not Allowed"))
		return
	}

	path := c.Param("id")
	if path == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.Writer.Write([]byte("Not Found"))
		return
	}

	if store[path] == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.Writer.Write([]byte("Not Found"))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, store[path])
	//c.Header("Location", store[path])
	//c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

//func CustomMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
//		res := httptest.NewRecorder()
//		next.ServeHTTP(res, r)
//
//		if res.Code == http.StatusMethodNotAllowed {
//			rw.WriteHeader(http.StatusMethodNotAllowed)
//			rw.Write([]byte("Method Not Allowed"))
//			return
//		}
//
//		for k, v := range res.Header() {
//			rw.Header()[k] = v
//		}
//
//		rw.WriteHeader(res.Code)
//
//		rw.Write(res.Body.Bytes())
//	})
//}
//
//func MainPage(rw http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	bodyBytes, err := io.ReadAll(r.Body)
//	defer r.Body.Close()
//	if err != nil {
//		http.Error(rw, "Bad Request", http.StatusBadRequest)
//		return
//	}
//
//	url := string(bodyBytes)
//	key := util.GetKey()
//
//	mtx.Lock()
//	store[key] = url
//	mtx.Unlock()
//
//	rw.Header().Set("Content-Type", "text/plain")
//	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(*config.BaseURL+"/"+key)))
//	rw.WriteHeader(http.StatusCreated)
//	rw.Write([]byte(*config.BaseURL + "/" + key))
//}
//
//func IDPage(rw http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	path := strings.TrimPrefix(r.URL.Path, "/")
//	if path == "" {
//		http.Error(rw, "Not Found", http.StatusNotFound)
//		return
//	}
//
//	if store[path] == "" {
//		http.Error(rw, "Not Found", http.StatusNotFound)
//		return
//	}
//
//	rw.Header().Set("Location", store[path])
//	rw.WriteHeader(http.StatusTemporaryRedirect)
//}
