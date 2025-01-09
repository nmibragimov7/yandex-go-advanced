package storage

import (
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/storage/file"
	"yandex-go-advanced/internal/storage/memory"

	"go.uber.org/zap"
)

type Storage interface {
	Get(key string) (string, error)
	Set(record *models.ShortenRecord) error
	Close() error
}

type StorageProvider struct {
	Config *config.Config
	Sugar  *zap.SugaredLogger
}

const (
	logKeyError = "error"
)

func (p *StorageProvider) CreateStorage() Storage {
	var str Storage

	if *p.Config.FilePath != "" {
		str, err := file.Init(*p.Config.FilePath)
		if err != nil {
			p.Sugar.Errorw(
				"failed to init file storage",
				logKeyError, err.Error(),
			)
		}

		return str
	}

	str = memory.Init()

	return str
}
