package main

import (
	"net/http"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/storage"
)

func main() {
	cnf := config.Init().GetConfig()
	sgr := logger.InitLogger()
	//gzp := &middleware.GzipProvider{}
	//lgp := &middleware.LoggerProvider{}
	str, err := storage.NewFileStorage(*cnf.FilePath)
	if err != nil {
		sgr.Errorw(
			"",
			"error", err.Error(),
		)
	}
	hdp := &handlers.HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
	}

	defer func() {
		err := str.Close()
		if err != nil {
			sgr.Errorw(
				"",
				"error", err.Error(),
			)
		}
	}()

	rtr := router.Provider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		//GzipMiddleware:   gzp,
		//LoggerMiddleWare: lgp,
		Handler: hdp,
	}

	sgr.Error(http.ListenAndServe(*cnf.Server, rtr.Router()))
}
