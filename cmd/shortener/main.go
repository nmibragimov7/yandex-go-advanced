package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/shutdown"
	"yandex-go-advanced/internal/storage"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

// logKeyError - error constant
// TimeForShutdown - timeout constant
const (
	logKeyError     = "error"
	TimeForShutdown = 3
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
			logKeyError, err.Error(),
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
	server := &http.Server{
		Addr:    *cnf.Server,
		Handler: rtr,
	}

	if cnf.HTTPS != nil && *cnf.HTTPS {
		certFile := "./cert.pem"
		keyFile := "./key.pem"

		if _, err = os.Stat(certFile); os.IsNotExist(err) {
			sgr.Errorw("HTTPS enabled but cert.pem not found", logKeyError, err.Error())
			return errors.New("cert.pem not found")
		}
		if _, err = os.Stat(keyFile); os.IsNotExist(err) {
			sgr.Errorw("HTTPS enabled but key.pem not found", logKeyError, err.Error())
			return errors.New("key.pem not found")
		}

		err = server.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			sgr.Errorw("failed to start server in HTTPS", logKeyError, err.Error())
			return errors.New("failed to start server in HTTPS")
		}
	}

	err = server.ListenAndServe()
	if err != nil {
		sgr.Errorw("failed to start server in HTTP", logKeyError, err.Error())
		return errors.New("failed to start server in HTTP")
	}

	shutdown.Shutdown(server, TimeForShutdown*time.Second)

	return nil
}
