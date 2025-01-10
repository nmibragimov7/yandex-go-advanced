package main

import (
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
)

const (
	logKeyError = "error"
)

func main() {
	cnf := config.Init()
	sgr := logger.Init()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			"error", err.Error(),
		)
	}
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &handlers.HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}
	rtr := router.RouterProvider{
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
	}

	sgr.Error(http.ListenAndServe(*cnf.Server, rtr.Router()))
}
