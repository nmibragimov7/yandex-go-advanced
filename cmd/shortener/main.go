package main

import (
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/storage"
)

func main() {
	server := config.Init()
	store := storage.NewStore()

	log.Fatal(http.ListenAndServe(server, handlers.Router(store)))
}
