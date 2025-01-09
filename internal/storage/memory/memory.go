package memory

import (
	"sync"
	"yandex-go-advanced/internal/models"
)

type Storage struct {
	storage map[string]string
	mtx     *sync.Mutex
}

func (s *Storage) Get(key string) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.storage[key], nil
}

func (s *Storage) Set(record *models.ShortenRecord) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.storage[record.ShortURL] = record.OriginalURL
	return nil
}

func (s *Storage) Close() error { return nil }

func Init() *Storage {
	return &Storage{
		storage: make(map[string]string),
		mtx:     &sync.Mutex{},
	}
}
