package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/router"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
)

func Example() {
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
		err = str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	hdp := &HandlerProvider{
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

	ts := httptest.NewServer(rtr.Router())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://practicum.yandex.ru/"))

	res, err := ts.Client().Do(req)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("failed to close body: %s", err.Error())
		}
	}()
	fmt.Println(res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}
