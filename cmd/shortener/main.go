package main

import (
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/middleware"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
)

func main() {
	cnf := config.Init().GetConfig()
	str := storage.NewStore()
	sgr := logger.InitLogger()
	mp := &middleware.Provider{}
	hp := &handlers.Provider{}

	log.Fatal(http.ListenAndServe(*cnf.Server, router.Router(cnf, str, sgr, mp, hp)))
}
