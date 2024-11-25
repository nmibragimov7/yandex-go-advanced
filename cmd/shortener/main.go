package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"yandex-go-advanced/internal/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", handlers.MainPage)
	r.Get(`/{id}`, handlers.IDPage)
	//mux := http.NewServeMux()
	//mux.HandleFunc(`/{id}`, handlers.IDPage)
	//mux.HandleFunc(`/`, handlers.MainPage)

	log.Fatal(http.ListenAndServe(":8080", r))
}
