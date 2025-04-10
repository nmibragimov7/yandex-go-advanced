package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"yandex-go-advanced/internal/storage/db/shortener"

	"yandex-go-advanced/internal/models"
)

type File interface {
	Seek(offset int64, whence int) (int64, error)
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	//io.Seeker
	//io.Reader
	//io.Writer
	//io.Closer
}

// Storage - struct that contains the necessary settings
type Storage struct {
	file   File
	writer *bufio.Writer
}

// Get - func for return record
func (s *Storage) Get(key string) (interface{}, error) {
	if _, err := s.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}

		if record.ShortURL == key {
			return &record, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return nil, errors.New("failed to find record")
}

// GetAll - func for return records
func (s *Storage) GetAll(_ interface{}) ([]interface{}, error) {
	return nil, nil
}

// Set - func for saving record in file
func (s *Storage) Set(record interface{}) (interface{}, error) {
	rec, ok := record.(*models.ShortenRecord)
	if !ok {
		return nil, errors.New("failed to parse record interface")
	}

	if _, err := s.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		var item models.ShortenRecord
		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}

		if item.OriginalURL == rec.OriginalURL {
			return nil, fmt.Errorf("shortener already exists: %w", shortener.NewDuplicateError(
				rec.ShortURL,
				"23505",
				errors.New("shortener already exists"),
			))
		}
	}

	data, err := json.Marshal(rec)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record: %w", err)
	}
	n, err := s.file.Write(append(data, '\n'))
	if err != nil {
		return nil, fmt.Errorf("failed to write record to file: %w", err)
	}

	return n, nil
}

// SetAll - func for saving records in file
func (s *Storage) SetAll(records []interface{}) error {
	rcs := make([]*models.ShortenRecord, 0, len(records))
	for _, record := range records {
		rec, ok := record.(*models.ShortenRecord)
		if !ok {
			return errors.New("failed to parse record interface")
		}

		rcs = append(rcs, rec)
	}

	for _, value := range rcs {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %w", err)
		}
		_, err = s.file.Write(append(data, '\n'))
		if err != nil {
			return fmt.Errorf("failed to write record to file: %w", err)
		}
	}

	return nil
}

// Close - func for close file
func (s *Storage) Close() error {
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

// Ping - func for ping file
func (s *Storage) Ping(_ context.Context) error { return nil }

// AddToChannel - func for add value in channel
func (s *Storage) AddToChannel(_ chan struct{}, _ ...chan interface{}) {}

var osOpenFile = os.OpenFile

// Init - initialize file instance
func Init(path string) (*Storage, error) {
	file, err := osOpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var record models.ShortenRecord
		if err = json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file record: %w", err)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner encountered an error: %w", err)
	}

	return &Storage{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}
