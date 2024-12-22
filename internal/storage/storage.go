package storage

import (
	"sync"
)

type Storage struct {
	storage map[string]string
	mtx     *sync.Mutex
}

func (s *Storage) Get() map[string]string {
	return s.storage
}

func (s *Storage) Save(key, url string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.storage[key] = url
}

func (s *Storage) GetByKey(key string) string {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.storage[key]
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
		mtx:     &sync.Mutex{},
	}
}
