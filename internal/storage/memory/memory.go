package memory

import (
	"context"
	"errors"
	"sync"
	internalModels "yandex-go-advanced/internal/models"
)

type Storage struct {
	storage map[string]string
	mtx     *sync.Mutex
}

func (s *Storage) Get(key string) (interface{}, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	var record models.ShortenRecord
	if s.storage[key] != "" {
		record = models.ShortenRecord{
			ShortURL:    key,
			OriginalURL: s.storage[key],
		}

		return &record, nil
	}

	return nil, errors.New("failed to find record")
}

func (s *Storage) GetAll(_ interface{}) ([]interface{}, error) {
	return nil, nil
}

func (s *Storage) Set(record interface{}) (interface{}, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	rec, ok := record.(*internalModels.ShortenRecord)
	if !ok {
		return nil, errors.New("failed to parse record interface")
	}

	s.storage[rec.ShortURL] = rec.OriginalURL

	return rec.OriginalURL, nil
}

func (s *Storage) SetAll(records []interface{}) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	rcs := make([]*internalModels.ShortenRecord, 0, len(records))
	for _, record := range records {
		rec, ok := record.(*internalModels.ShortenRecord)
		if !ok {
			return errors.New("failed to parse record interface")
		}

		rcs = append(rcs, rec)
	}

	for _, value := range rcs {
		s.storage[value.ShortURL] = value.OriginalURL
	}

	return nil
}

func (s *Storage) Close() error { return nil }

func (s *Storage) Ping(_ context.Context) error { return nil }

func (s *Storage) AddToChannel(_ chan struct{}, _ ...chan interface{}) {}

func Init() *Storage {
	return &Storage{
		storage: make(map[string]string),
		mtx:     &sync.Mutex{},
	}
}
