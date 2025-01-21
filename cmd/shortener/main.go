package main

import (
	"fmt"
	"log"
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
)

const (
	logKeyError = "error"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cnf := config.Init()
	sgr := logger.Init()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			"error", err.Error(),
		)

		return fmt.Errorf("failed to init storage: %w", err)
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
	ssp := &session.SessionProvider{
		Config: cnf,
	}
	rtr := router.RouterProvider{
		Storage: str,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	sgr.Error(http.ListenAndServe(*cnf.Server, rtr.Router()))

	return nil
}
