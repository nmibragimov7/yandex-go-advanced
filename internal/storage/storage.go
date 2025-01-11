package storage

import (
	"context"
	"fmt"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/storage/db"
	"yandex-go-advanced/internal/storage/file"
	"yandex-go-advanced/internal/storage/memory"
)

type Storage interface {
	Get(entity string, key string) (interface{}, error)
	Set(entity string, record interface{}) error
	SetByTransaction(entity string, records []interface{}) error
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
		value, err := storage.Get(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get record from database: %w", err)
		}

		return value, nil
	}

	if p.file != nil {
		value, err := p.file.Get(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get record from file: %w", err)
		}

		return value, nil
	}

	value, err := p.memory.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get record from file: %w", err)
	}

	return value, nil
}

func (p *StorageProvider) Set(entity string, record interface{}) error {
	if storage, ok := p.db[entity]; ok {
		err := storage.Set(record)
		if err != nil {
			return fmt.Errorf("failed to save record to database: %w", err)
		}

		return nil
	}

	if p.file != nil {
		err := p.file.Set(record)
		if err != nil {
			return fmt.Errorf("failed to save record to file: %w", err)
		}

		return nil
	}

	err := p.memory.Set(record)
	if err != nil {
		return fmt.Errorf("failed to save record to memory: %w", err)
	}

	return nil
}

func (p *StorageProvider) SetByTransaction(entity string, records []interface{}) error {
	if storage, ok := p.db[entity]; ok {
		err := storage.SetByTransaction(records)
		if err != nil {
			return fmt.Errorf("failed to save records to database: %w", err)
		}

		return nil
	}

	if p.file != nil {
		err := p.file.SetByTransaction(records)
		if err != nil {
			return fmt.Errorf("failed to save records to file: %w", err)
		}

		return nil
	}

	err := p.memory.SetByTransaction(records)
	if err != nil {
		return fmt.Errorf("failed to save records to memory: %w", err)
	}

	return nil
}

func (p *StorageProvider) Close() error {
	if p.file != nil {
		if err := p.file.Close(); err != nil {
			return fmt.Errorf("failed to close file storage: %w", err)
		}
	}

	for entity, storage := range p.db {
		if err := storage.Close(); err != nil {
			return fmt.Errorf("failed to close database for %s: %w", entity, err)
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

	return nil
}

func Init(cnf *config.Config) (Storage, error) {
	memoryStorage := memory.Init()

	var fileStorage *file.Storage
	if *cnf.FilePath != "" {
		var err error
		fileStorage, err = file.Init(*cnf.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize file storage: %w", err)
		}
	}

	dbStorages := make(map[string]*db.Storage)
	if *cnf.DataBase != "" {
		database, err := db.Init(*cnf.DataBase)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}

		if database != nil {
			err := db.Bootstrap(database.DB)
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
