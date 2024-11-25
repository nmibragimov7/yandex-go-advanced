package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
)

func main() {
	config.Init()
	flag.Parse()

	r := chi.NewRouter()

	r.Post("/", handlers.MainPage)
	r.Get(`/{id}`, handlers.IDPage)

	//mux := http.NewServeMux()
	//mux.HandleFunc(`/{id}`, handlers.IDPage)
	//mux.HandleFunc(`/`, handlers.MainPage)

	fmt.Println("Server", *config.Server)

	log.Fatal(http.ListenAndServe(*config.Server, r))
}
