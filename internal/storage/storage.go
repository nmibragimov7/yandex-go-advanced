package storage

import (
	"context"
	"errors"
	"fmt"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/storage/db"
	"yandex-go-advanced/internal/storage/file"
	"yandex-go-advanced/internal/storage/memory"
)

type Storage interface {
	Get(entity string, key string) (interface{}, error)
	Set(entity string, record interface{}) error
	Close() error
	Ping(ctx context.Context) error
}

type StorageProvider struct {
	db     map[string]*db.Storage
	file   *file.Storage
	memory *memory.Storage
}

func (p *StorageProvider) Get(entity string, key string) (interface{}, error) {
	if storage, ok := p.db[entity]; ok {
		return storage.Get(key)
	}

	if p.file != nil {
		return p.file.Get(key)
	}

	return p.memory.Get(key)
}

func (p *StorageProvider) Set(entity string, record interface{}) error {
	if storage, ok := p.db[entity]; ok {
		return storage.Set(record)
	}

	if p.file != nil {
		return p.file.Set(record)
	}

	return p.memory.Set(record)
}

func (p *StorageProvider) Close() error {
	if p.file != nil {
		if err := p.file.Close(); err != nil {
			return fmt.Errorf("failed to close file storage: %w", err)
		}
	}

	for entity, storage := range p.db {
		if err := storage.Close(); err != nil {
			return fmt.Errorf("failed to close db storage for %s: %w", entity, err)
		}
	}

	return nil
}

func (p *StorageProvider) Ping(ctx context.Context) error {
	for entity, storage := range p.db {
		if err := storage.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping db storage for %s: %w", entity, err)
		}
	}

	return errors.New("failed to ping database")
}

func Init(config *config.Config) (Storage, error) {
	memoryStorage := memory.Init()

	var fileStorage *file.Storage
	if *config.FilePath != "" {
		var err error
		fileStorage, err = file.Init(*config.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize file storage: %w", err)
		}
	}

	dbStorages := make(map[string]*db.Storage)
	if *config.DataBase != "" {
		database, err := db.Init(*config.DataBase)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}

		if database != nil {
			err := db.InitTables(database.DB)
			if err != nil {
				return nil, fmt.Errorf("failed to create table queries: %w", err)
			}

			dbStorages["shortener"] = &db.Storage{DB: database.DB}
		}
	}

	storage := &StorageProvider{
		db:     dbStorages,
		file:   fileStorage,
		memory: memoryStorage,
	}
	return storage, nil
}
