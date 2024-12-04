package storage

import "sync"

type Store struct {
	Store map[string]string
	mtx   *sync.Mutex
}

func (s *Store) GetStore() map[string]string {
	return s.Store
}
func (s *Store) SaveStore(key, url string) {
	s.mtx.Lock()
	s.Store[key] = url
	defer s.mtx.Unlock()
}

func NewStore() *Store {
	return &Store{
		Store: make(map[string]string),
		mtx:   &sync.Mutex{},
	}
}
