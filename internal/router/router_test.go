package router

import (
	"testing"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/session"

	"github.com/stretchr/testify/assert"
)

func ptr(s string) *string {
	return &s
}

func TestRouter(t *testing.T) {
	t.Run("Router initialized", func(t *testing.T) {
		cnf := &config.Config{
			Server:   ptr(":8080"),
			BaseURL:  ptr("http://localhost:8080"),
			FilePath: ptr("./storage.txt"),
			DataBase: ptr("host=localhost user=postgres password=admin dbname=postgres sslmode=disable"),
		}
		sgr := logger.Init()
		ssp := &session.SessionProvider{
			Config: cnf,
		}
		hdp := &handlers.HandlerProvider{
			Config:  cnf,
			Storage: nil,
			Sugar:   sgr,
			Session: ssp,
		}
		rtp := RouterProvider{
			Storage: nil,
			Config:  cnf,
			Sugar:   sgr,
			Handler: hdp,
			Session: ssp,
		}

		router := rtp.Router()
		assert.NotNil(t, router)
	})
}
