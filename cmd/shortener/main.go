package main

import (
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/storage"
)

func main() {
	cnf := config.Init().GetConfig()
	str := storage.NewStore()
	sgr := logger.InitLogger()

	log.Fatal(http.ListenAndServe(*cnf.Server, handlers.Router(cnf, str, sgr)))
}
