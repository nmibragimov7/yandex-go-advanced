package main

import (
	"fmt"
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
)

func main() {
	server := config.Init()
	fmt.Println("server", server)

	log.Fatal(http.ListenAndServe(server, handlers.Router()))
}
