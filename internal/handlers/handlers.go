package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/util"
)

var store = make(map[string]string)
var mtx sync.Mutex

func CustomMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		res := httptest.NewRecorder()
		next.ServeHTTP(res, r)

		if res.Code == http.StatusMethodNotAllowed {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Method Not Allowed"))
			return
		}

		for k, v := range res.Header() {
			rw.Header()[k] = v
		}

		rw.WriteHeader(res.Code)

		rw.Write(res.Body.Bytes())
	})
}

func Router() chi.Router {
	r := chi.NewRouter()
	r.Use(CustomMiddleware)

	r.Post("/", MainPage)
	r.Get(`/{id}`, IDPage)

	return r
}

func MainPage(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	url := string(bodyBytes)
	key := util.GetKey()

	mtx.Lock()
	store[key] = url
	mtx.Unlock()

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(*config.BaseURL+"/"+key)))
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(*config.BaseURL + "/" + key))
}

func IDPage(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}

	if store[path] == "" {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Location", store[path])
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
