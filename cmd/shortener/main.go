package main

import (
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/db"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
)

const (
	logKeyError = "error"
)

func main() {
	cnf := config.Init().GetConfig()
	sgr := logger.InitLogger()
	str, err := storage.InitFileStorage(*cnf.FilePath)
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()
	if err != nil {
		sgr.Errorw(
			"failed to init file storage",
			logKeyError, err.Error(),
		)
	}

	dbp := db.DatabaseProvider{
		Sugar:  sgr,
		Config: cnf,
	}
	err = dbp.Init()
	if err != nil {
		sgr.Errorw(
			"failed to init database",
			logKeyError, err.Error(),
		)
		return
	}
	database := dbp.Get()
	defer func() {
		err := database.Close()
		if err != nil {
			sgr.Errorw(
				"Failed to close database connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &handlers.HandlerProvider{
		Config:   cnf,
		Storage:  str,
		Sugar:    sgr,
		Database: &dbp,
	}
	rtr := router.RouterProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		Handler: hdp,
	}

	sgr.Error(http.ListenAndServe(*cnf.Server, rtr.Router()))
}
