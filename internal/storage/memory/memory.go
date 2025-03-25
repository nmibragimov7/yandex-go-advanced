package memory

import (
	"context"
	"errors"
	"sync"

	"yandex-go-advanced/internal/models"
)

// Storage - struct that contains the necessary settings
type Storage struct {
	storage map[string]string
	mtx     *sync.Mutex
}

// Get - func for return record
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

// GetAll - func for return records
func (s *Storage) GetAll(_ interface{}) ([]interface{}, error) {
	return nil, nil
}

// Set - func for saving record in memory
func (s *Storage) Set(record interface{}) (interface{}, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	rec, ok := record.(*models.ShortenRecord)
	if !ok {
		return nil, errors.New("failed to parse record interface")
	}

	s.storage[rec.ShortURL] = rec.OriginalURL

	return rec.OriginalURL, nil
}

// SetAll - func for saving records in memory
func (s *Storage) SetAll(records []interface{}) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	rcs := make([]*models.ShortenRecord, 0, len(records))
	for _, record := range records {
		rec, ok := record.(*models.ShortenRecord)
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

// Close - func for close memory
func (s *Storage) Close() error { return nil }

// Ping - func for ping memory
func (s *Storage) Ping(_ context.Context) error { return nil }

// AddToChannel - func for add value in channel
func (s *Storage) AddToChannel(_ chan struct{}, _ ...chan interface{}) {}

// Init - initialize memory instance
func Init() *Storage {
	return &Storage{
		storage: make(map[string]string),
		mtx:     &sync.Mutex{},
	}
}
