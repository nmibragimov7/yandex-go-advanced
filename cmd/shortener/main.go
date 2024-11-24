package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var store = make(map[string]string)
var mtx sync.Mutex

func getKey() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)[:8]
}

func mainPage(rw http.ResponseWriter, r *http.Request) {
	baseUrl := "http://localhost:8080/"

	if r.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(rw, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	//	return
	//}

	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	url := string(bodyBytes)
	key := getKey()

	mtx.Lock()
	store[key] = url
	mtx.Unlock()

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(baseUrl+key)))
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(baseUrl + key))
}

func idPage(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(rw, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	//	return
	//}

	path := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println("path", path)
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, idPage)
	mux.HandleFunc(`/`, mainPage)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
