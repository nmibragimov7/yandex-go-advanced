package main

import (
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/storage/db"
)

const (
	logKeyError = "error"
)

func main() {
	cnf := config.Init()
	sgr := logger.Init()

	dbp := db.DatabaseProvider{
		Config: cnf,
		Sugar:  sgr,
	}
	database, err := dbp.Init()
	if err != nil {
		sgr.Errorw(
			"failed to init database",
			logKeyError, err.Error(),
		)
	}
	defer func() {
		if database != nil {
			err := database.Close()
			if err != nil {
				sgr.Errorw(
					"Failed to close database connection",
					logKeyError, err.Error(),
				)
			}
		}
	}()
	if database != nil {
		err := dbp.CreateTables(database)
		if err != nil {
			sgr.Errorw(
				"failed to create table query",
				"error", err.Error(),
			)
		}
	}

	stp := storage.StorageProvider{
		Config: cnf,
		Sugar:  sgr,
	}

	str := stp.CreateStorage()
	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage",
				"error", err.Error(),
			)
		}
	}()

	hdp := &handlers.HandlerProvider{
		Config:   cnf,
		Storage:  str,
		Sugar:    sgr,
		Database: database,
	}
	rtr := router.RouterProvider{
		Config:  cnf,
		Sugar:   sgr,
		Handler: hdp,
	}

	sgr.Error(http.ListenAndServe(*cnf.Server, rtr.Router()))
}
