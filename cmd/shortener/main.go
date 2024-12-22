package main

import (
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
	sgr := logger.InitLogger()
	mp := &middleware.Provider{}
	hp := &handlers.Provider{}
	str, err := storage.NewFileStorage(*cnf.FilePath)

	sugar := sgr.Get()
	if err != nil {
		sugar.Errorw(
			"",
			"error", err.Error(),
		)
	}
	defer func() {
		err := str.Close()
		if err != nil {
			sugar.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()

	sugar.Fatal(http.ListenAndServe(*cnf.Server, router.Router(cnf, str, sgr, mp, hp)))
}
