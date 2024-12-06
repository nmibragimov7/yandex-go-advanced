package storage

import (
	"sync"
)

type Store struct {
	Store map[string]string
	mtx   *sync.Mutex
}

func (s *Store) SaveStore(key, url string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.Store[key] = url
}

func (s *Store) Get(key string) string {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.Store[key]
}

func NewStore() *Store {
	return &Store{
		Store: make(map[string]string),
		mtx:   &sync.Mutex{},
	}
}
