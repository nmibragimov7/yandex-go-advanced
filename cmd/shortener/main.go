package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

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
		err = str.Close()
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
	sgr.Log(1, "Build version: ", buildVersion)
	sgr.Log(1, "Build date: ", buildDate)
	sgr.Log(1, "Build commit: ", buildCommit)

	rtr := rtp.Router()

	if cnf.HTTPS != nil && *cnf.HTTPS {
		certFile := "../../cert.pem"
		keyFile := "../../key.pem"

		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			sgr.Errorw("HTTPS enabled but cert.pem not found", logKeyError, err.Error())
			return fmt.Errorf("cert.pem not found")
		}
		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			sgr.Errorw("HTTPS enabled but key.pem not found", logKeyError, err.Error())
			return fmt.Errorf("key.pem not found")
		}

		sgr.Error(http.ListenAndServeTLS(*cnf.Server, certFile, keyFile, rtr))
	}
	sgr.Error(http.ListenAndServe(*cnf.Server, rtr))

	return nil
}
