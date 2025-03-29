package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

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
	defer func() {
		err := sgr.Sync()
		if err != nil {
			log.Printf("failed to sync logger: %s", err.Error())
		}
	}()

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

	ssp := &session.SessionProvider{
		Config: cnf,
	}
	hdp := &handlers.HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		Session: ssp,
	}
	rtp := router.RouterProvider{
		Storage: str,
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
		Session: ssp,
	}

	sgr.Log(1, "server started in: ", *cnf.Server)
	sgr.Error(http.ListenAndServe(*cnf.Server, rtp.Router()))

	return nil
}
