package pkg

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
)

func Config() {
	config.Init()
	flag.Parse()
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", handlers.MainPage)
	r.Get(`/{id}`, handlers.IDPage)

	return r
}
