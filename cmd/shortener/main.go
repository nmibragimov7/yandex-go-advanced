package main

import (
	"fmt"
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/pkg"
)

func main() {
	//mux := http.NewServeMux()x
	//mux.HandleFunc(`/{id}`, handlers.IDPage)
	//mux.HandleFunc(`/`, handlers.MainPage)

	pkg.ParseFlag()
	fmt.Println("Server", *config.Server)

	log.Fatal(http.ListenAndServe(*config.Server, handlers.Router()))
}
